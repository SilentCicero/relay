package unittest

import (
	"bytes"
	"fmt"
	"gitlab.com/jonas.jasas/httprelay/pkg/model"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func newProxyCliData() (proxyCliData *model.ProxyCliData, r *http.Request, data []byte, serId, serPath, query, fragment, url string) {
	data = bytes.Repeat([]byte{10}, 10000)
	serId = "123"
	serPath = "/test"
	query = "first=1&second=last"
	fragment = "frag1=first&frag2=second"
	url = fmt.Sprintf("https://demo.httprelay.io/proxy/%s%s?%s#%s", serId, serPath, query, fragment)
	r, _ = http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	proxyCliData = model.NewProxyCliData(r, serId, serPath)
	return
}

func TestNewProxyCliData(t *testing.T) {
	pcd, r, data, serId, serPath, query, fragment, url := newProxyCliData()

	if pcd.Method != http.MethodPost {
		t.Error("HTTP method mismatch")
	}

	if pcd.Scheme != "https" {
		t.Error("Scheme mismatch")
	}

	if pcd.Host != "demo.httprelay.io" {
		t.Error("Host mismatch")
	}

	if pcd.Url != url {
		t.Error("URL mismatch")
	}

	if pcd.Path != serPath {
		t.Error("Path mismatch")
	}

	if pcd.Query != query {
		t.Error("Query mismatch")
	}

	if pcd.Fragment != fragment {
		t.Error("Fragment mismatch")
	}

	if pcd.SerId != serId {
		t.Error("SerId mismatch")
	}

	if pcd.Header != &r.Header {
		t.Error("Header mismatch")
	}

	if b, err := ioutil.ReadAll(pcd.Body); err == nil {
		if bytes.Compare(b, data) != 0 {
			t.Error("Body data mismatch")
		}
	} else {
		t.Error("Error while reading body")
	}
}

func TestProxyCliDataClose(t *testing.T) {
	pcd, _, _, _, _, _, _, _ := newProxyCliData()
	if !pcd.CloseRespChan() {
		t.Fail()
	}
	if pcd.CloseRespChan() {
		t.Fail()
	}

	pcd, _, _, _, _, _, _, _ = newProxyCliData()
	go func() { pcd.RespChan <- nil }()
	time.Sleep(100000)
	if !pcd.CloseRespChan() {
		t.Fail()
	}
	if pcd.CloseRespChan() {
		t.Fail()
	}
}
