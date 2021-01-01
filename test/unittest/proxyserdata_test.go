package unittest

import (
	"bytes"
	"gitlab.com/jonas.jasas/httprelay/pkg/model"
	"io/ioutil"
	"net/http"
	"testing"
)

func newProxySerData() (proxyCliData *model.ProxySerData, r *http.Request, data []byte) {
	data = bytes.Repeat([]byte{10}, 10000)
	r, _ = http.NewRequest(http.MethodPost, "https://domain/proxy/123/test", bytes.NewReader(data))
	proxyCliData = model.NewProxySerData(r)
	return
}

func TestNewProxySerData(t *testing.T) {
	psd, r, data := newProxySerData()

	if psd.Header != &r.Header {
		t.Fail()
	}

	if b, err := ioutil.ReadAll(psd.Body); err == nil {
		if bytes.Compare(b, data) != 0 {
			t.Fail()
		}
	} else {
		t.Fail()
	}
}
