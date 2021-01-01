package model

import "net/http"

type LinkData struct {
	Data     *PtpData
	BackChan chan *Meta
}

func NewLinkData(r *http.Request) *LinkData {
	return &LinkData{
		Data:     NewPtpData(r),
		BackChan: make(chan *Meta),
	}
}
