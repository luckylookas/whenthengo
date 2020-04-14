package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestYamlParser_Parse(t *testing.T) {
	parser := YamlParser{}
	file, err := os.Open("test_resources/" + t.Name() + ".yml")
	assert.NoError(t, err)
	defer file.Close()
	actual, err := parser.Parse(file)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(actual))

	assert.Equal(t,  "/path/test", actual[0].When.URL)
	assert.Equal(t, 2, len(actual[0].When.Headers))
	assert.Equal(t,"application/json", actual[0].When.Headers["Accept"][0])
	//json parses as float not int
	assert.Equal(t, "6", actual[0].When.Headers["Content-Length"][0])
	assert.Equal(t, "abc\ndef\n", actual[0].When.Body)

	assert.Equal(t, 200, actual[0].Then.Status)
	assert.Equal(t, 100, actual[0].Then.Delay)
	assert.Equal(t, "k", actual[0].Then.Body)
	assert.Equal(t, 1, len(actual[0].Then.Headers))
	assert.Equal(t, "1", actual[0].Then.Headers["Content-Length"][0])
}

func TestYamlParser_Parse_HeaderCasting(t *testing.T) {
	parser := YamlParser{}
	file, err := os.Open("test_resources/" + t.Name() + ".yml")
	assert.NoError(t, err)
	defer file.Close()
	actual, err := parser.Parse(file)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(actual))

	assert.Equal(t, actual[0].When.Headers["String"], []string{"json"})
	assert.Equal(t, actual[0].When.Headers["Int"], []string{"1"})
	assert.Equal(t, actual[0].When.Headers["Float"], []string{"1.23"})
	assert.Equal(t, actual[0].When.Headers["IntSlice"], []string{"1", "2"})
	assert.Equal(t, actual[0].When.Headers["StringSlice"], []string{"json", "yaml"})


	assert.Equal(t, actual[0].Then.Headers["String"], []string{"json"})
	assert.Equal(t, actual[0].Then.Headers["Int"], []string{"1"})
	assert.Equal(t, actual[0].Then.Headers["Float"], []string{"1.23"})
	assert.Equal(t, actual[0].Then.Headers["IntSlice"], []string{"1", "2"})
	assert.Equal(t, actual[0].Then.Headers["StringSlice"], []string{"json", "yaml"})
}

