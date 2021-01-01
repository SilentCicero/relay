package repository

import (
	"gitlab.com/jonas.jasas/httprelay/pkg/model"
	"sync"
)

type syncMap map[string]chan *model.SyncData

type SyncRep struct {
	syncMap
	sync.RWMutex
	model.Waiters
	stopChan <-chan struct{}
}

func NewSyncRep(stopChan <-chan struct{}) *SyncRep {
	return &SyncRep{
		syncMap{},
		sync.RWMutex{},
		model.Waiters{},
		stopChan,
	}
}

func (sr *SyncRep) Conduct(id string, syncData *model.SyncData, cancelChan <-chan struct{}) (ptpData *model.PtpData, ok bool) {
	sr.AddWaiter()
	defer sr.RemoveWaiter()

	sr.Lock()
	syncDataChan, exist := sr.syncMap[id]
	if exist {
		// B - peer
		delete(sr.syncMap, id)
	} else {
		// A - peer
		syncDataChan = make(chan *model.SyncData, 1)
		syncDataChan <- syncData
		sr.syncMap[id] = syncDataChan
	}
	sr.Unlock()

	defer sr.closeAndRemoveSync(id, syncDataChan)

	if exist {
		// B - peer
		select {
		case peerSyncData := <-syncDataChan:
			// As chan is buffered so one value always should be present
			ptpData = peerSyncData.Data
			select {
			case peerSyncData.BackChan <- syncData.Data:
				ok = true
			case <-syncDataChan: // A peer exited
			case <-cancelChan:
			case <-sr.stopChan:
			}
		default:
			// Something is horribly wrong, as one value always should be present
		}
	} else {
		// A - peer
		select {
		case ptpData = <-syncData.BackChan:
			ok = true
		case <-cancelChan:
		case <-sr.stopChan:
		}
	}

	return
}

func (sr *SyncRep) Count() int {
	sr.RLock()
	defer sr.RUnlock()
	return len(sr.syncMap)
}

func (sr *SyncRep) closeAndRemoveSync(id string, c chan *model.SyncData) {
	select {
	case <-c:
	default:
		close(c)
	}

	sr.Lock()
	defer sr.Unlock()
	delete(sr.syncMap, id)
}
