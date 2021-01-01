package testlib

import (
	"gitlab.com/jonas.jasas/httprelay/pkg/controller"
	"gitlab.com/jonas.jasas/httprelay/pkg/repository"
	"io"
	"net/http"
	"net/http/httptest"
)

func NewProxyCtrl() (proxyCtrl *controller.ProxyCtrl, stopChan chan struct{}, closeChan chan struct{}) {
	stopChan = make(chan struct{})
	proxyRep := repository.NewProxyRep()
	proxyCtrl = controller.NewProxyCtrl(proxyRep, stopChan)
	closeChan = make(chan struct{})
	return
}

func ProxyCtrlCliReq(ctrl *controller.ProxyCtrl, url string, header map[string]string, dataReader io.Reader) *httptest.ResponseRecorder {
	r := newReq(http.MethodPost, url, header, dataReader)
	defer r.Body.Close()
	w := httptest.NewRecorder()
	ctrl.Conduct(w, r)
	return w
}

func ProxyCtrlSerReq(ctrl *controller.ProxyCtrl, url string, header map[string]string, dataReader io.Reader, reqJobId string, wSecret string) (resp *httptest.ResponseRecorder, respJobId string) {
	r := newReq("SERVE", url, header, dataReader)
	defer r.Body.Close()
	r.Header.Add("HttpRelay-Proxy-JobId", reqJobId)
	r.Header.Add("HttpRelay-WSecret", wSecret)
	w := httptest.NewRecorder()
	ctrl.Conduct(w, r)
	respJobId = w.Header().Get("HttpRelay-Proxy-JobId")
	return w, respJobId
}
