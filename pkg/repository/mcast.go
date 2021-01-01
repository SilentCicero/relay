package repository

import (
	"gitlab.com/jonas.jasas/httprelay/pkg/model"
	"sync"
)

type McastSeqMap map[string]*model.McastSeq

type McastRep struct {
	McastSeqMap
	sync.RWMutex
	stopChan <-chan struct{}
}

func NewMcastRep(stopChan <-chan struct{}) *McastRep {
	result := &McastRep{
		McastSeqMap: McastSeqMap{},
		stopChan:    stopChan,
	}
	return result
}

func (mr *McastRep) GetData(id string, seqId int) (data *model.TeeData, ok bool) {
	mr.RLock()
	mcasts, ok := mr.McastSeqMap[id]
	mr.RUnlock()

	if ok {
		data, ok = mcasts.GetData(seqId)
	}

	return
}

func (mr *McastRep) getOrCreateMcasts(id string) (mcasts *model.McastSeq) {
	mcasts, ok := mr.McastSeqMap[id]
	if !ok {
		initialSeqId := 0

		mcasts = model.NewMcastSeq(initialSeqId)
		mr.McastSeqMap[id] = mcasts
	}

	return
}

func (mr *McastRep) Write(id string, data *model.TeeData, wSecret string) (seqId int, ok bool) {
	mr.Lock()
	defer mr.Unlock()
	mcastSeq := mr.getOrCreateMcasts(id)
	if ok = mcastSeq.WAuth(wSecret); ok {
		seqId = mcastSeq.Write(data)
	}
	return
}

func (mr *McastRep) Read(id string, wantedSeqId int, cancelChan <-chan struct{}) (data *model.TeeData, seqId int, ok bool) {
	mr.Lock()
	mcasts := mr.getOrCreateMcasts(id)
	mr.Unlock()
	return mcasts.Read(wantedSeqId, cancelChan)
}

func (mr *McastRep) removeOutdated() {
	mr.Lock()
	defer mr.Unlock()
	for k, v := range mr.McastSeqMap {
		if v.Expired() {
			delete(mr.McastSeqMap, k)
			v.Close()
		}
	}
}

func (mr *McastRep) Size() (size int) {
	mr.RLock()
	defer mr.RUnlock()

	for _, v := range mr.McastSeqMap {
		size += v.Size()
	}
	return
}

func (mr *McastRep) DataCount() (size int) {
	mr.RLock()
	defer mr.RUnlock()

	for _, v := range mr.McastSeqMap {
		size += v.DataCount()
	}
	return
}

func (mr *McastRep) WaiterCount() (size int) {
	mr.RLock()
	defer mr.RUnlock()

	for _, v := range mr.McastSeqMap {
		size += v.WaiterCount()
	}
	return
}
