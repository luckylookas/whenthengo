package main

import (
	"os"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestJsonParser_Parse(t *testing.T) {
	parser := JsonParser{}
	file, err := os.Open("test_resources/" + t.Name() + ".json")
	assert.NoError(t, err)
	defer file.Close()
	actual, err := parser.Parse(file)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(actual))
	assert.Equal(t, "/path/test", actual[0].When.URL)
	assert.Equal(t, 2, len(actual[0].When.Headers))
	assert.Equal(t,"application/json", actual[0].When.Headers["Accept"])
	//json parses as float not int
	assert.Equal(t, "6", actual[0].When.Headers["Content-Length"])
	assert.Equal(t, "abc\ndef\n", actual[0].When.Body)

	assert.Equal(t, 200, actual[0].Then.Status)
	assert.Equal(t, 100, actual[0].Then.Delay)
	assert.Equal(t, "k", actual[0].Then.Body)
	assert.Equal(t, 1, len(actual[0].Then.Headers))
	assert.Equal(t, "1", actual[0].Then.Headers["Content-Length"])
}