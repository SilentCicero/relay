package model

import (
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

type Meta struct {
	Time          time.Time
	Header        http.Header
	ContentType   string
	Method        string
	Query         string
	SrcIP         string
	SrcPort       string
	ContentLength int64
	// After altering structure Size() method must be updated
}

func NewMeta(r *http.Request) *Meta {
	t := time.Now()
	contentType := r.Header.Get("Content-Type")
	srcIP := r.Header.Get("X-Real-IP")
	srcPort := r.Header.Get("X-Real-Port")
	if srcIP == "" {
		srcIP, srcPort, _ = net.SplitHostPort(r.RemoteAddr)
	}

	query := filterQuery(r.URL.Query())
	method := r.Method

	return &Meta{t, r.Header, contentType, method, query.Encode(), srcIP, srcPort, r.ContentLength}
}

func filterQuery(query url.Values) (filtered url.Values) {
	filtered = url.Values{}
	for k, vals := range query {
		if !strings.EqualFold(k, "wsecret") && !strings.EqualFold(k, "seqid") {
			for _, v := range vals {
				filtered.Add(k, v)
			}
		}
	}
	return
}

func (m *Meta) Memory() int64 {
	structSize := uint64(unsafe.Sizeof(m))
	stringSize := len(m.ContentType) + len(m.Method) + len(m.Query) + len(m.SrcIP) + len(m.SrcPort)
	return int64(structSize) + int64(stringSize)
}

func (m *Meta) WriteHeaders(w http.ResponseWriter, yourTime time.Time, content bool) {
	if content {
		if m.ContentType != "" {
			w.Header().Set("Content-Type", m.ContentType)
		}
		if m.ContentLength > -1 {
			// If length is known is better to set it or "Transfer-Encoding: chunked" will be used
			w.Header().Set("Content-Length", strconv.FormatInt(m.ContentLength, 10))
		}
	}
	w.Header().Set("X-Real-IP", m.SrcIP)
	w.Header().Set("X-Real-Port", m.SrcPort)
	w.Header().Set("HttpRelay-Time", toUnixMilli(m.Time))
	w.Header().Set("HttpRelay-Your-Time", toUnixMilli(yourTime))
	w.Header().Set("HttpRelay-Method", m.Method)
	if m.Query != "" {
		w.Header().Set("HttpRelay-Query", m.Query)
	}
}

func toUnixMilli(t time.Time) string {
	mills := t.UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
	return strconv.FormatInt(mills, 10)
}
