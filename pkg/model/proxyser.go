package model

import "sync"

// Creating when first request to/from server is made
type ProxySer struct {
	ReqChan chan *ProxyCliData
	jobMap  map[string]*ProxyCliData
	jobMapL sync.Mutex
	comm
}

func NewProxySer() *ProxySer {
	return &ProxySer{
		ReqChan: make(chan *ProxyCliData),
		jobMap:  map[string]*ProxyCliData{},
		comm:    newComm(),
	}
}

func (ps *ProxySer) AddJob(jobId string, proxyData *ProxyCliData) {
	ps.jobMapL.Lock()
	defer ps.jobMapL.Unlock()
	ps.jobMap[jobId] = proxyData
}

func (ps *ProxySer) TakeJob(jobId string) (proxyData *ProxyCliData, ok bool) {
	ps.jobMapL.Lock()
	defer ps.jobMapL.Unlock()
	if proxyData, ok = ps.jobMap[jobId]; ok {
		delete(ps.jobMap, jobId)
	}
	return
}

func (ps *ProxySer) RemoveJob(cliData *ProxyCliData) {
	ps.jobMapL.Lock()
	defer ps.jobMapL.Unlock()

	for k, v := range ps.jobMap {
		if v == cliData {
			delete(ps.jobMap, k)
			break
		}
	}
}
