package controller

import (
	"gitlab.com/jonas.jasas/httprelay/pkg/model"
	"gitlab.com/jonas.jasas/httprelay/pkg/repository"
	"net/http"
	"strings"
)

type ProxyCtrl struct {
	rep      *repository.ProxyRep
	stopChan <-chan struct{}
	*model.Waiters
}

func NewProxyCtrl(rep *repository.ProxyRep, stopChan <-chan struct{}) *ProxyCtrl {
	return &ProxyCtrl{
		rep:      rep,
		stopChan: stopChan,
		Waiters:  model.NewWaiters(),
	}
}

func (pc *ProxyCtrl) Conduct(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	select {
	case <-pc.stopChan:
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	default:
	}

	serId, serPath := pc.parsePath(r.URL.Path)
	ser := pc.rep.GetSer(serId)

	ser.AddWaiter()
	defer ser.RemoveWaiter()

	if strings.EqualFold(r.Method, "SERVE") {
		pc.handleServer(ser, r, w)
	} else {
		pc.handleClient(r, w, ser, serId, serPath)
	}
}

func (pc *ProxyCtrl) parsePath(path string) (serId, serPath string) {
	p := strings.TrimLeft(path, "/")
	arr := strings.SplitN(p, "/", 3)
	serId = arr[1]
	serPath = "/"
	if len(arr) > 2 {
		serPath += arr[2]
	}
	return
}
