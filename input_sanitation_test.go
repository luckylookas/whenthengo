package main

import (
	"bytes"
	"gotest.tools/assert"
	"testing"
)

func TestString_CleanUrl(t *testing.T) {
	assert.Equal(t, "abc/d", cleanUrl("/abc/d/"))
	assert.Equal(t, "abc/d", cleanUrl("//abc/d// /"))
	assert.Equal(t, "abc", cleanUrl(" /abc "))
}

func TestString_CleanMethod(t *testing.T) {
	assert.Equal(t, "get", cleanUrl("GET"))
	assert.Equal(t, "get", cleanUrl("get"))
}

func TestString_CleanBodyString(t *testing.T) {
	assert.Equal(t, cleanBodyString(`
		{
			"some": "json",  
				  "with":  "spaces"
		}`), `{"some":"json","with":"spaces"}`)
}


func TestString_CleanBodyBytes(t *testing.T) {
	assert.Equal(t,
		0,
		bytes.Compare(
		cleanBodyBytes([]byte(`
		{
			"some": "json",  
				  "with":  "spaces"
		}`)), []byte(`{"some":"json","with":"spaces"}`)))
}

func TestString_CleanHeaderKey(t *testing.T) {
	assert.Equal(t, cleanHeaderValue("Content-Type"), "content-type")
}

func TestString_CleanHeaderValue(t *testing.T) {
	assert.Equal(t, cleanHeaderValue(`Application/Json; encoding="UTF-8"`), `application/jsonencoding="utf-8"`)
}

func TestString_CleanHeaders(t *testing.T) {
	headers := map[string][]string{
		"Accept": {"APPLICATION/JSON"},
		"Stuff":  {"1;2,3 4"},
	}

	headers = cleanHeaders(headers)

	assert.Equal(t, "application/json", headers["accept"][0])
	assert.Equal(t, "1234", headers["stuff"][0])
}

