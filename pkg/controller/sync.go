package controller

import (
	"gitlab.com/jonas.jasas/httprelay/pkg/model"
	"net/http"
	"strings"
	"time"
)

type SyncRep interface {
	Conduct(id string, syncData *model.SyncData, cancelChan <-chan struct{}) (ptpData *model.PtpData, ok bool)
}

type SyncCtrl struct {
	rep      SyncRep
	stopChan <-chan struct{}
	*model.Waiters
}

func NewSyncCtrl(rep SyncRep, stopChan <-chan struct{}) *SyncCtrl {
	return &SyncCtrl{
		rep:      rep,
		stopChan: stopChan,
		Waiters:  model.NewWaiters(),
	}
}

func (sc *SyncCtrl) Conduct(w http.ResponseWriter, r *http.Request) {
	yourTime := time.Now()

	sc.AddWaiter()
	defer sc.RemoveWaiter()
	defer r.Body.Close()

	select {
	case <-sc.stopChan:
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	default:
	}

	pathArr := strings.Split(r.URL.Path, "/")
	linkId := pathArr[len(pathArr)-1]

	syncData := model.NewSyncData(r)
	if ptpData, ok := sc.rep.Conduct(linkId, syncData, r.Context().Done()); ok {
		<-syncData.Data.Content.Buff()
		ptpData.Write(w, yourTime)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
}
