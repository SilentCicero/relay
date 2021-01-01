package integration

import (
	"bytes"
	"fmt"
	"gitlab.com/jonas.jasas/closechan"
	"gitlab.com/jonas.jasas/httprelay/pkg/controller"
	"gitlab.com/jonas.jasas/httprelay/test/testlib"
	"gitlab.com/jonas.jasas/rwmock"
	"math/rand"
	"net/http/httptest"
	"strconv"
	"sync"
	"testing"
	"time"
)

type proxyTestData [][]byte

func TestProxy(t *testing.T) {
	servers := []string{"first", "second", "third", "fourth", "fifth"}
	proxyData := genProxyData()
	ctrl, _, _ := testlib.NewProxyCtrl()
	closeChan := closechan.NewCloseChan()
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			runProxyCliReq(t, ctrl, servers, proxyData, closeChan)
			wg.Done()
		}()
	}

	for i := 0; i < 5; i++ {
		for _, server := range servers {
			go func(ser string) {
				runProxySerReq(t, ctrl, ser, proxyData, closeChan)
				closeChan.Close()
			}(server)
		}
	}

	wg.Wait()
	closeChan.Close()
}

func genProxyData() (data proxyTestData) {
	data = make(proxyTestData, 100)
	for i := 0; i < len(data); i++ {
		b := make([]byte, rand.Intn(1000000))
		rand.Read(b)
		data[i] = b
	}
	return
}

func runProxyCliReq(t *testing.T, proxyCtrl *controller.ProxyCtrl, servers []string, data proxyTestData, closeChan *closechan.CloseChan) {
	for i := 0; i < 20; i++ {
		for _, ser := range servers {
			w, path := newProxyCliReq(proxyCtrl, data, ser)
			if closeChan.Closed() {
				return
			}
			if w.Header().Get("test-path") != path {
				t.Error("client received incorrect path in response")
				return
			}
			if dataIdx, err := strconv.Atoi(w.Header().Get("test-data-idx")); err == nil {
				if !testlib.RespDataEq(w.Body, data[dataIdx]) {
					t.Error("client received incorrect body in response")
					return
				}
			} else {
				t.Error("client received incorrect data array index")
				return
			}
		}
	}
}

func runProxySerReq(t *testing.T, proxyCtrl *controller.ProxyCtrl, server string, data proxyTestData, closeChan *closechan.CloseChan) {
	path, jobId := "", ""
	var w *httptest.ResponseRecorder
	for !closeChan.Closed() {
		w, jobId = newProxySerReq(proxyCtrl, data, server, path, jobId)
		path = w.Header().Get("test-path")
		jobId = w.Header().Get("HttpRelay-Proxy-JobId")

		if dataIdx, err := strconv.Atoi(w.Header().Get("test-data-idx")); err == nil {
			if !testlib.RespDataEq(w.Body, data[dataIdx]) {
				t.Error("server received incorrect body in response")
				return
			}
		} else {
			t.Error("server received incorrect data array index")
			return
		}
	}
}

func newProxyCliReq(proxyCtrl *controller.ProxyCtrl, data proxyTestData, ser string) (w *httptest.ResponseRecorder, path string) {
	path = genId(20)
	dataIdx := rand.Intn(len(data) - 1)
	url := fmt.Sprintf("https://example.com/proxy/%s/%s", ser, path)
	header := map[string]string{
		"test-path":     path,
		"test-data-idx": strconv.Itoa(dataIdx),
	}

	r := rwmock.NewShaperRand(bytes.NewReader(data[dataIdx]), 1, len(data[dataIdx]), 0, time.Millisecond)
	w = testlib.ProxyCtrlCliReq(proxyCtrl, url, header, r)
	return
}

func newProxySerReq(proxyCtrl *controller.ProxyCtrl, data proxyTestData, ser, path, jobId string) (w *httptest.ResponseRecorder, respJobId string) {
	dataIdx := rand.Intn(len(data) - 1)
	url := fmt.Sprintf("https://example.com/proxy/%s", ser)
	header := map[string]string{
		"HttpRelay-Proxy-Headers": "test-data-idx, test-path",
		"test-path":               path,
		"test-data-idx":           strconv.Itoa(dataIdx),
	}

	r := rwmock.NewShaperRand(bytes.NewReader(data[dataIdx]), 1, len(data[dataIdx]), 0, time.Millisecond)
	w, respJobId = testlib.ProxyCtrlSerReq(proxyCtrl, url, header, r, jobId, "")
	return
}
