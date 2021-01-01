package controller

import (
	"fmt"
	"gitlab.com/jonas.jasas/httprelay/pkg/model"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type McastRep interface {
	Read(id string, wantedSeqId int, cancelChan <-chan struct{}) (data *model.TeeData, seqId int, ok bool)
	Write(id string, data *model.TeeData, wSecret string) (seqId int, ok bool)
}

type McastCtrl struct {
	rep      McastRep
	stopChan <-chan struct{}
	*model.Waiters
}

func NewMcastCtrl(rep McastRep, stopChan <-chan struct{}) *McastCtrl {
	return &McastCtrl{
		rep:      rep,
		stopChan: stopChan,
		Waiters:  model.NewWaiters(),
	}
}

const seqIdParamName = "SeqId"

func (mc *McastCtrl) Conduct(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	pathArr := strings.Split(r.URL.Path, "/")
	mcastId := pathArr[len(pathArr)-1]

	closeChan := r.Context().Done()

	yourTime := time.Now()
	if strings.EqualFold(r.Method, http.MethodGet) {
		seqId := seqId(r)
		cache := seqId > -1

		if data, seqId, ok := mc.rep.Read(mcastId, seqId, closeChan); ok {
			writeHeaders(w, r, seqId, data, yourTime, cache)
			io.Copy(w, data.Content.NewReader())
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
		}
		return
	} else if strings.EqualFold(r.Method, http.MethodPost) {
		data := model.NewTeeData(r)
		if seqId, ok := mc.rep.Write(mcastId, data, wSecret(r)); ok {
			if _, err := data.CopyContent(); err == nil {
				w.Header().Set("HttpRelay-SeqId", strconv.Itoa(seqId))
				data.Meta.WriteHeaders(w, yourTime, false)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}
		return
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func writeHeaders(w http.ResponseWriter, r *http.Request, seqId int, data *model.TeeData, yourTime time.Time, cache bool) {
	if cache {
		w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d", model.MaxAge/time.Second))
		w.Header().Set("Vary", "Cookie")
	}

	nextSeqid := seqId + 1
	w.Header().Set("Set-Cookie", fmt.Sprintf("%s=%d; Path=%s; SameSite=None; Secure", seqIdParamName, nextSeqid, r.URL.Path))
	w.Header().Set("HttpRelay-SeqId", strconv.Itoa(seqId))
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	data.Meta.WriteHeaders(w, yourTime, true)
}

func seqId(r *http.Request) (seqId int) {
	seqId = -1

	for param := range r.URL.Query() {
		if strings.EqualFold(param, seqIdParamName) {
			seqId, _ = strconv.Atoi(r.URL.Query().Get(param))
			return
		}
	}

	if cookie, err := r.Cookie(seqIdParamName); err == nil {
		seqId, _ = strconv.Atoi(cookie.Value)
	}
	return
}
