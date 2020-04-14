package main

import (
	"net/http"
)

type mockResponseWriter struct {
	header http.Header
	status *int
	Body   []byte
}

func (m mockResponseWriter) Header() http.Header {
	return m.header
}

func (m mockResponseWriter) Write(a []byte) (int, error) {
	copy(m.Body, a)
	return len(a), nil
}

func (m mockResponseWriter) WriteHeader(statusCode int) {
	*m.status = statusCode
}

//todo