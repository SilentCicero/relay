package integration

import (
	"bytes"
	"fmt"
	"gitlab.com/jonas.jasas/httprelay/pkg/controller"
	"gitlab.com/jonas.jasas/httprelay/pkg/repository"
	"math/rand"
	"net/http"
	"testing"
)

type syncReqData struct {
	AData []byte
	BData []byte
}

func NewSyncData(n int) map[string]*syncReqData {
	syncMap := make(map[string]*syncReqData)
	for i := 0; i < n; i++ {
		d := syncReqData{
			make([]byte, rand.Intn(100000)),
			make([]byte, rand.Intn(100000)),
		}
		rand.Read(d.AData)
		rand.Read(d.BData)
		syncMap[genId(10)] = &d
	}
	return syncMap
}

func TestSync(t *testing.T) {
	syncReqDataMap := NewSyncData(10000)

	cancelChan := make(chan struct{})
	syncRep := repository.NewSyncRep(cancelChan)
	syncCtrl := controller.NewSyncCtrl(syncRep, cancelChan)
	handler := http.HandlerFunc(syncCtrl.Conduct)

	resChan := make(chan bool)
	for k, v := range syncReqDataMap {
		go func(k string, v *syncReqData) {
			resA, resB := doReqPair(handler, "POST", "POST", fmt.Sprintf("/sync/%s", k), fmt.Sprintf("/sync/%s", k), v.AData, v.BData)
			resChan <- bytes.Compare(resA, v.BData) == 0 && bytes.Compare(resB, v.AData) == 0
		}(k, v)
	}

	for i := 0; i < len(syncReqDataMap); i++ {
		if !<-resChan {
			t.Fail()
		}
	}

	if syncRep.Count() != 0 {
		t.Fail()
	}
}
