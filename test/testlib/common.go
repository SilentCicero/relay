package testlib

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
)

func newReq(method string, url string, header map[string]string, dataReader io.Reader) (r *http.Request) {
	r, _ = http.NewRequest(method, url, dataReader)
	if header != nil {
		for k, v := range header {
			r.Header.Add(k, v)
		}
	}
	return
}

func RespDataEq(body io.Reader, data []byte) bool {
	respData, _ := ioutil.ReadAll(body)
	return bytes.Compare(respData, data) == 0
}
