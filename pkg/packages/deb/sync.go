package deb

import (
	"pault.ag/go/archive"
	"sync"
)

// SyncOptions provides an option set to manage a sync group of multiple producer workers and
// single consumer, leveraging Go sync.WaitGroup and channels to notify errors, results and completion
// of consuming the results from the single consumer worker.
type SyncOptions struct {
	waitGroup *sync.WaitGroup
	resultCh  chan []archive.Package
	errCh     chan error
	doneCh    chan bool
}

func NewSyncOptions(waitGroup *sync.WaitGroup, packagesCh chan []archive.Package, errCh chan error, doneCh chan bool) *SyncOptions {
	return &SyncOptions{
		waitGroup: waitGroup,
		resultCh:  packagesCh,
		errCh:     errCh,
		doneCh:    doneCh,
	}
}

func (o *SyncOptions) WaitGroup() *sync.WaitGroup {
	return o.waitGroup
}

func (o *SyncOptions) ResultCh() chan []archive.Package {
	return o.resultCh
}

func (o *SyncOptions) ErrCh() chan error {
	return o.errCh
}

func (o *SyncOptions) DoneCh() chan bool {
	return o.doneCh
}

func (o *SyncOptions) WaitAll() {

	// Wait for producers to complete.
	o.waitGroup.Wait()
	close(o.resultCh)
	close(o.errCh)

	// Wait for consumers to complete.
	<-o.doneCh
}
