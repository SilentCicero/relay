package model

import (
	"gitlab.com/jonas.jasas/buffreader"
	"io"
	"net/http"
	"time"
)

type PtpData struct {
	Meta    *Meta
	Content *buffreader.BuffReader
}

func NewPtpData(r *http.Request) *PtpData {
	return &PtpData{
		Meta:    NewMeta(r),
		Content: buffreader.New(r.Body),
	}
}

func (pd *PtpData) Write(w http.ResponseWriter, yourTime time.Time) (err error) {
	pd.Meta.WriteHeaders(w, yourTime, true)
	_, err = io.Copy(w, pd.Content)
	return
}
