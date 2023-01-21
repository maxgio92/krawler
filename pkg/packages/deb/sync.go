package deb

import (
	"pault.ag/go/archive"
	"sync"
)

// MPSCQueue provides an option set to manage a sync group of multiple producer workers and
// single consumer, leveraging Go sync.WaitGroup and channels to notify errors, results and completion
// of consuming the results from the single consumer worker.
type MPSCQueue struct {
	producersWG    *sync.WaitGroup
	consumerDoneCh chan bool
	msgCh          chan []archive.Package
	errCh          chan error
}

// func NewSyncOptions(waitGroup *sync.WaitGroup, packagesCh chan []archive.Package, errCh chan error, doneCh chan bool) *MPSCQueue {
func NewSyncOptions(parallelism int) *MPSCQueue {
	wg := &sync.WaitGroup{}
	wg.Add(parallelism)

	msgCh := make(chan []archive.Package)

	errCh := make(chan error)

	doneCh := make(chan bool, 1)

	return &MPSCQueue{
		producersWG:    wg,
		consumerDoneCh: doneCh,
		msgCh:          msgCh,
		errCh:          errCh,
	}
}

// SendMessage sends a message as variadic parameter msg of type archive.Package to the messages queue.
func (q *MPSCQueue) SendMessage(msg ...archive.Package) {
	q.msgCh <- msg
}

// SendMessageAndComplete sends a message as variadic parameter msg of type archive.Package to the messages queue,
// and eventually signals the completion of the current producer.
func (q *MPSCQueue) SendMessageAndComplete(msg ...archive.Package) {
	defer q.SigProducerCompletion()
	q.msgCh <- msg
}

// SendError sends an error message of type error to the errors queue.
func (q *MPSCQueue) SendError(err error) {
	q.errCh <- err
}

// Consume listens for both messages and errors on queues and do something with them,
// as specified by msgHandler and errHandler functions.
func (q *MPSCQueue) Consume(msgHandler func(msg ...archive.Package), errHandler func(err error)) {
	for q.errCh != nil || q.msgCh != nil {
		select {
		case p, ok := <-q.msgCh:

			// If the channel is still open.
			if ok {

				// Do something with the message.
				msgHandler(p...)
				continue
			}
			q.msgCh = nil
		case e, ok := <-q.errCh:

			// If the channel is still open.
			if ok {

				// Do something with error.
				errHandler(e)
				continue
			}
			q.errCh = nil
		}
	}
	q.SigConsumerCompletion()
}

// SigProducerCompletion signals that a producer completed its work.
func (q *MPSCQueue) SigProducerCompletion() {
	q.producersWG.Done()
}

// SigConsumerCompletion signals that the consumer completed its work.
func (q *MPSCQueue) SigConsumerCompletion() {
	q.consumerDoneCh <- true
}

func (q *MPSCQueue) WaitAndClose() {

	// Wait for producersWG to complete.
	q.producersWG.Wait()
	close(q.msgCh)
	close(q.errCh)

	// Wait for consumers to complete.
	<-q.consumerDoneCh
}

func (q *MPSCQueue) ProducersWG() *sync.WaitGroup {
	return q.producersWG
}

func (q *MPSCQueue) MessageCh() chan []archive.Package {
	return q.msgCh
}

func (q *MPSCQueue) ErrorCh() chan error {
	return q.errCh
}

func (q *MPSCQueue) ConsumerDoneCh() chan bool {
	return q.consumerDoneCh
}
