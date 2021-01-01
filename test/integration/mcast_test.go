package integration

import (
	"bytes"
	"fmt"
	"gitlab.com/jonas.jasas/httprelay/pkg/controller"
	"gitlab.com/jonas.jasas/httprelay/pkg/repository"
	"gitlab.com/jonas.jasas/rwmock"
	"io/ioutil"
	"math/rand"
	"net/http"
	"sync"
	"testing"
	"time"
)

type mcastTestData map[string][][]byte

func genMcastData() (data mcastTestData) {
	data = make(mcastTestData)

	const cnt = 100
	for i := 0; i < cnt; i++ {
		seqLen := rand.Intn(100)
		seq := make([][]byte, seqLen)

		for s := 0; s < seqLen; s++ {
			b := make([]byte, rand.Intn(100000))
			rand.Read(b)
			seq = append(seq, b)
		}

		data[genId(10)] = seq
	}
	return
}

func TestMcast(t *testing.T) {
	cancelChan := make(chan struct{})
	mcastRep := repository.NewMcastRep(cancelChan)
	mcastCtrl := controller.NewMcastCtrl(mcastRep, cancelChan)
	handler := http.HandlerFunc(mcastCtrl.Conduct)

	data := genMcastData()
	wg := sync.WaitGroup{}
	for k, v := range data {
		wg.Add(1)
		go func(id string, seq [][]byte) {
			for _, b := range seq {
				time.Sleep(time.Duration(rand.Int63n(1000000)))
				r := rwmock.NewShaperRand(bytes.NewReader(b), 1, 1000, 0, time.Millisecond)
				doReq(handler, "POST", fmt.Sprintf("/mcast/%s", id), r)
			}
			wg.Done()
		}(k, v)

		for rc := 0; rc < 10; rc++ {
			wg.Add(1)
			go func(id string, seq [][]byte) {
				time.Sleep(time.Duration(rand.Int63n(1000000000)))
				for n, b := range seq {
					resp := doReq(handler, "GET", fmt.Sprintf("/mcast/%s?SeqId=%d", id, n), nil)
					r := rwmock.NewShaperRand(resp.Body, 1, 1000, 0, time.Millisecond)
					resBuf, _ := ioutil.ReadAll(r)
					if bytes.Compare(resBuf, b) != 0 {
						t.Error("Read data mismatch")
					}
				}
				wg.Done()
			}(k, v)
		}
	}

	wg.Wait()
	fmt.Print(mcastRep.Size(), mcastRep.DataCount())
}
