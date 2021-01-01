package model

import (
	"gitlab.com/jonas.jasas/bufftee"
	"io"
	"net/http"
)

type TeeData struct {
	r          *http.Request
	Meta       *Meta
	Content    *bufftee.BuffTee
	cancelChan chan struct{}
}

func NewTeeData(r *http.Request) *TeeData {
	cancelChan := make(chan struct{})
	return &TeeData{
		r:          r,
		Meta:       NewMeta(r),
		Content:    bufftee.NewBuffTee(cancelChan),
		cancelChan: cancelChan,
	}
}

func (td *TeeData) CopyContent() (n int64, err error) {
	n, err = io.Copy(td.Content, td.r.Body)
	if err == nil {
		td.Content.Close()
	} else {
		td.Cancel()
	}
	return
}

func (td *TeeData) Cancel() {
	select {
	case <-td.cancelChan:
	default:
		close(td.cancelChan)
	}
}
