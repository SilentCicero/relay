package unittest

import (
	"bytes"
	"gitlab.com/jonas.jasas/httprelay/test/testlib"
	"net/http"
	"testing"
)

const serUrl = "https://domain/proxy/123"
const serPath = "/test"

var cliData1 = newData("Client data 1. ", 10000)
var cliData2 = newData("Client data 2. ", 10000)
var serData1 = newData("Server data 1. ", 10000)
var serData2 = newData("Server data 2. ", 10000)

func TestProxyCtrlConduct(t *testing.T) {
	proxyCtrl, _, closeChan := testlib.NewProxyCtrl()
	go func() {
		defer close(closeChan)
		resp := testlib.ProxyCtrlCliReq(proxyCtrl, serUrl+serPath, nil, bytes.NewReader(cliData1))
		if !testlib.RespDataEq(resp.Body, serData1) {
			t.Error("Client received wrong response body")
		}
	}()
	go func() {
		defer close(closeChan)
		resp, jobId := testlib.ProxyCtrlSerReq(proxyCtrl, serUrl, nil, bytes.NewReader([]byte{}), "", "")
		if resp.Header().Get("HttpRelay-Proxy-Path") != serPath {
			t.Error("Server received wrong path")
			return
		}
		if !testlib.RespDataEq(resp.Body, cliData1) {
			t.Error("Server received wrong client data body")
			return
		}
		testlib.ProxyCtrlSerReq(proxyCtrl, serUrl, nil, bytes.NewReader(serData1), jobId, "")
	}()
	<-closeChan
}

func TestProxyCtrlWSecret(t *testing.T) {
	proxyCtrl, _, closeChan := testlib.NewProxyCtrl()
	go func() {
		defer close(closeChan)
		resp := testlib.ProxyCtrlCliReq(proxyCtrl, serUrl, nil, bytes.NewReader(cliData1))
		if !testlib.RespDataEq(resp.Body, serData1) {
			t.Error("Client received wrong response body")
			return
		}
		resp = testlib.ProxyCtrlCliReq(proxyCtrl, serUrl, nil, bytes.NewReader(cliData2))
		if !testlib.RespDataEq(resp.Body, serData2) {
			t.Error("Client received wrong response body")
			return
		}
	}()

	go func() {
		const goodSecret = "secret1"
		const badSecret = "bad secret"
		defer close(closeChan)
		resp, jobId := testlib.ProxyCtrlSerReq(proxyCtrl, serUrl, nil, bytes.NewReader([]byte{}), "", goodSecret)
		if !testlib.RespDataEq(resp.Body, cliData1) {
			t.Error("Server received wrong client data body")
			return
		}

		resp, _ = testlib.ProxyCtrlSerReq(proxyCtrl, serUrl, nil, bytes.NewReader(serData1), jobId, badSecret)
		if resp.Code != http.StatusUnauthorized {
			t.Error("Server is accessing unauthorized data channel with the bad secret")
			return
		}

		resp, jobId = testlib.ProxyCtrlSerReq(proxyCtrl, serUrl, nil, bytes.NewReader(serData1), jobId, goodSecret)
		if !testlib.RespDataEq(resp.Body, cliData2) {
			t.Error("Server received wrong client data body")
			return
		}

		resp, _ = testlib.ProxyCtrlSerReq(proxyCtrl, serUrl, nil, bytes.NewReader(serData1), jobId, badSecret)
		if resp.Code != http.StatusUnauthorized {
			t.Error("Server is accessing unauthorized data channel with the bad secret")
			return
		}

		resp, jobId = testlib.ProxyCtrlSerReq(proxyCtrl, serUrl, nil, bytes.NewReader(serData2), jobId, goodSecret)
	}()
	<-closeChan
}

func TestProxyCtrlHeaders(t *testing.T) {
	proxyCtrl, _, closeChan := testlib.NewProxyCtrl()

	cliHeader := map[string]string{
		"test-client-header1": "Should pass 1",
		"test-client-header2": "Should pass 2",
	}

	serHeaderPass := map[string]string{
		"test-server-header1": "Should pass 1",
		"test-server-header3": "Should pass 1",
	}

	serHeaderNoPass := map[string]string{
		"test-server-header2":     "Should not pass",
		"HttpRelay-Proxy-Headers": "test-server-header1, test-server-header3",
	}

	serHeader := map[string]string{}
	for k, v := range serHeaderPass {
		serHeader[k] = v
	}
	for k, v := range serHeaderNoPass {
		serHeader[k] = v
	}

	go func() {
		defer close(closeChan)
		resp := testlib.ProxyCtrlCliReq(proxyCtrl, serUrl, cliHeader, bytes.NewReader(cliData1))
		if !testlib.RespDataEq(resp.Body, serData1) {
			t.Error("Client received wrong response body")
		}
		for k, v := range serHeaderPass {
			if resp.Header().Get(k) != v {
				t.Error("Client didn't receive expected server headers")
				return
			}
		}
		for k, _ := range serHeaderNoPass {
			if resp.Header().Get(k) != "" {
				t.Error("Client received unexpected server header")
				return
			}
		}
	}()

	go func() {
		defer close(closeChan)
		resp, jobId := testlib.ProxyCtrlSerReq(proxyCtrl, serUrl, nil, bytes.NewReader([]byte{}), "", "")
		if !testlib.RespDataEq(resp.Body, cliData1) {
			t.Error("Server received wrong client data body")
			return
		}
		for k, v := range cliHeader {
			if resp.Header().Get(k) != v {
				t.Error("Server didn't receive expected client header")
				return
			}
		}
		testlib.ProxyCtrlSerReq(proxyCtrl, serUrl, serHeader, bytes.NewReader(serData1), jobId, "")
	}()
	<-closeChan
}
