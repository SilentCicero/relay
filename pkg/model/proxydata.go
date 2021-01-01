package model

import (
	"gitlab.com/jonas.jasas/buffreader"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

type ProxyCliData struct {
	Url       string
	Method    string
	Scheme    string
	Host      string
	Port      string
	Path      string
	Query     string
	Fragment  string
	SerId     string
	Header    *http.Header
	Body      *buffreader.BuffReader
	RespChan  chan *ProxySerData
	RespChanL sync.Mutex
}

func NewProxyCliData(r *http.Request, serId, serPath string) (proxyReqData *ProxyCliData) {
	protoHeader := r.Header.Get("X-Forwarded-Proto")
	scheme := protoHeader
	if scheme == "" {
		scheme = r.URL.Scheme
	}
	if scheme == "" {
		scheme = "http"
	}
	scheme = strings.TrimSpace(scheme)

	hostArr := strings.Split(r.Host, ":")
	remoteAddrArr := strings.Split(r.RemoteAddr, ":")

	host := r.Header.Get("X-Forwarded-Host")
	if host == "" {
		host = hostArr[0]
	}
	if host == "" {
		host = remoteAddrArr[0]
	}
	host = strings.TrimSpace(host)

	port := r.Header.Get("X-Forwarded-Port")
	if port == "" && protoHeader == "http" {
		port = "80"
	}
	if port == "" && protoHeader == "https" {
		port = "443"
	}
	if port == "" && len(hostArr) == 2 {
		port = hostArr[1]
	}
	if port == "" && len(remoteAddrArr) == 2 {
		port = remoteAddrArr[1]
	}
	port = strings.TrimSpace(port)

	fullUrl, _ := url.Parse(r.URL.String()) // Making a copy of request URL
	if !fullUrl.IsAbs() {
		fullUrl.Scheme = scheme
		if port != "80" && port != "443" {
			fullUrl.Host = host + ":" + port
		} else {
			fullUrl.Host = host
		}
	}

	proxyReqData = &ProxyCliData{
		Url:      fullUrl.String(),
		Method:   r.Method,
		Scheme:   scheme,
		Host:     host,
		Port:     port,
		Path:     serPath,
		Query:    r.URL.RawQuery,
		Fragment: r.URL.Fragment,
		SerId:    serId,
		Header:   &r.Header,
		Body:     buffreader.New(r.Body),
	}

	proxyReqData.Body.Buff()
	proxyReqData.RespChan = make(chan *ProxySerData)

	return
}

func (pcd *ProxyCliData) CloseRespChan() (ok bool) {
	pcd.RespChanL.Lock()
	defer pcd.RespChanL.Unlock()

	select {
	case _, ok = <-pcd.RespChan:
	default:
		ok = true
	}

	if ok {
		close(pcd.RespChan)
	}

	return
}

type ProxySerData struct {
	Header *http.Header
	Body   *buffreader.BuffReader
}

func NewProxySerData(r *http.Request) (proxyRespData *ProxySerData) {
	proxyRespData = &ProxySerData{
		Header: &r.Header,
		Body:   buffreader.New(r.Body),
	}
	proxyRespData.Body.Buff()
	return
}
