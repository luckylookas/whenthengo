package types

import (
	"bytes"
	"gotest.tools/assert"
	"testing"
)

func TestString_CleanUrl(t *testing.T) {
	assert.Equal(t, "abc/d", CleanUrl("/abc/d/"))
	assert.Equal(t, "abc/d", CleanUrl("//abc/d// /"))
	assert.Equal(t, "abc", CleanUrl(" /abc "))
}

func TestString_CleanMethod(t *testing.T) {
	assert.Equal(t, "get", CleanMethod("GET"))
	assert.Equal(t, "get", CleanMethod("get"))
}

func TestString_CleanBodyString(t *testing.T) {
	assert.Equal(t, CleanBodyString(`
		{
			"some": "json",  
				  "with":  "spaces"
		}`), `{"some":"json","with":"spaces"}`)
}


func TestString_CleanBodyBytes(t *testing.T) {
	assert.Equal(t,
		0,
		bytes.Compare(
		CleanBodyBytes([]byte(`
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

	headers = CleanHeaders(headers)

	assert.Equal(t, "application/json", headers["accept"][0])
	assert.Equal(t, "1234", headers["stuff"][0])
}

