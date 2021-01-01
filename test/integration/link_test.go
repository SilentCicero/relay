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

func TestLink(t *testing.T) {
	genId(10)

	cancelChan := make(chan struct{})
	linkRep := repository.NewLinkRep(cancelChan)
	syncCtrl := controller.NewLinkCtrl(linkRep, cancelChan)
	handler := http.HandlerFunc(syncCtrl.Conduct)

	resChan := make(chan bool)
	const cnt = 10000
	for i := 0; i < cnt; i++ {
		go func() {
			id := genId(10)
			AData := make([]byte, rand.Intn(100000))
			_, resB := doReqPair(handler, "POST", "GET", fmt.Sprintf("/sync/%s", id), fmt.Sprintf("/sync/%s", id), AData, nil)
			resChan <- bytes.Compare(resB, AData) == 0
		}()
	}

	for i := 0; i < cnt; i++ {
		if !<-resChan {
			t.Fail()
		}
	}
}
