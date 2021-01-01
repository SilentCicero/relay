package integration

import (
	"bytes"
	"gitlab.com/jonas.jasas/rwmock"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"time"
)

func genId(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func doReq(handler http.HandlerFunc, method, path string, body io.Reader) *httptest.ResponseRecorder {
	resp := httptest.NewRecorder()
	req := httptest.NewRequest(method, path, body)

	handler.ServeHTTP(resp, req)

	return resp
}

func doReqPair(handler http.HandlerFunc, methodA, methodB, pathA, pathB string, bodyA, bodyB []byte) ([]byte, []byte) {

	respAChan := make(chan *httptest.ResponseRecorder, 1)

	rA := rwmock.NewShaperRand(bytes.NewReader(bodyA), 1, 1000, 0, time.Microsecond)
	rB := rwmock.NewShaperRand(bytes.NewReader(bodyB), 1, 1000, 0, time.Microsecond)

	go func() {
		time.Sleep(time.Duration(rand.Int63n(1000000)))
		respAChan <- doReq(handler, methodA, pathA, rA)
	}()

	time.Sleep(time.Duration(rand.Int63n(1000000)))
	respB := doReq(handler, methodB, pathB, rB)

	respA := <-respAChan

	return respA.Body.Bytes(), respB.Body.Bytes()
}
