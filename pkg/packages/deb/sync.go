package deb

import (
	"pault.ag/go/archive"
	"sync"
)

type SyncOptions struct {
	waitGroup  *sync.WaitGroup
	packagesCh chan []archive.Package
	errCh      chan error
}

func NewSyncOptions(waitGroup *sync.WaitGroup, packagesCh chan []archive.Package, errCh chan error) *SyncOptions {
	return &SyncOptions{
		waitGroup:  waitGroup,
		packagesCh: packagesCh,
		errCh:      errCh,
	}
}

func (o *SyncOptions) WaitGroup() *sync.WaitGroup {
	return o.waitGroup
}
