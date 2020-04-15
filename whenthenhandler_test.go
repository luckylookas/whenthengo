package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github/luckylukas/whenthengo/types"
	"io"
	"net/http"
	"testing"
)

func NewMockWriter(nested io.Writer) mockResponseWriter{
	return mockResponseWriter{
		header: make(http.Header),
		status: new(int),
		Body:  nested,
	}
}

type mockResponseWriter struct {
	header http.Header
	status *int
	Body   io.Writer
}

func (m mockResponseWriter) Header() http.Header {
	return m.header
}

func (m mockResponseWriter) Write(a []byte) (int, error) {
	return m.Body.Write(a)
}

func (m mockResponseWriter) WriteHeader(statusCode int) {
	*m.status = statusCode
}

func TestWriteThen(t *testing.T) {
	actual := bytes.Buffer{}
	writer := NewMockWriter(&actual)

	writeThen(writer, &types.Then{
		Status:  100,
		Delay:   200,
		Headers: types.Header{
			"some-data": []string{"some", "values"},
		},
		Body:    `{"content":"a"}`,
	})

	assert.Equal(t, 100, *writer.status)
	assert.Equal(t, "some", writer.header.Get("some-data"))
	assert.ElementsMatch(t, []string{"some", "values"}, writer.header["Some-Data"])
	assert.ElementsMatch(t, actual.Bytes(), []byte(`{"content":"a"}`))
}