package model

import (
	"net/http"
)

type SyncData struct {
	Data     *PtpData
	BackChan chan *PtpData
}

func NewSyncData(r *http.Request) *SyncData {
	return &SyncData{
		Data:     NewPtpData(r),
		BackChan: make(chan *PtpData),
	}
}
