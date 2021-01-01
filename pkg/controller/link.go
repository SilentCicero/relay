package controller

import (
	"gitlab.com/jonas.jasas/httprelay/pkg/model"
	"net/http"
	"strings"
	"time"
)

type LinkRep interface {
	Read(id string, r *http.Request, cancelChan <-chan struct{}) (ptpData *model.PtpData, ok bool)
	Write(id string, linkData *model.LinkData, wSecret string, cancelChan <-chan struct{}) (meta *model.Meta, ok bool, auth bool)
}

type LinkCtrl struct {
	rep      LinkRep
	stopChan <-chan struct{}
	*model.Waiters
}

func NewLinkCtrl(rep LinkRep, stopChan <-chan struct{}) *LinkCtrl {
	return &LinkCtrl{
		rep:      rep,
		stopChan: stopChan,
		Waiters:  model.NewWaiters(),
	}
}

func (lc *LinkCtrl) Conduct(w http.ResponseWriter, r *http.Request) {
	lc.AddWaiter()
	defer lc.RemoveWaiter()

	defer r.Body.Close()

	select {
	case <-lc.stopChan:
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	default:
	}

	pathArr := strings.Split(r.URL.Path, "/")
	id := pathArr[len(pathArr)-1]

	yourTime := time.Now()

	if strings.EqualFold(r.Method, http.MethodGet) {
		if ptpData, ok := lc.rep.Read(id, r, r.Context().Done()); ok {
			lc.AddWaiter()
			ptpData.Write(w, yourTime)
			lc.RemoveWaiter()
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
		}
	} else if strings.EqualFold(r.Method, http.MethodPost) {
		linkData := model.NewLinkData(r)
		if meta, ok, auth := lc.rep.Write(id, linkData, wSecret(r), r.Context().Done()); ok && auth {
			<-linkData.Data.Content.Buff()
			meta.WriteHeaders(w, yourTime, false)
		} else {
			if auth {
				w.WriteHeader(http.StatusServiceUnavailable)
			} else {
				w.WriteHeader(http.StatusUnauthorized)
			}
		}
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
