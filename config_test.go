package main

import (
	"gotest.tools/assert"
	"testing"
)

func TestLoadConfig(t *testing.T) {

	config := LoadConfig()

	assert.Equal(t, config.Port, "80")
	assert.Equal(t, config.WhenThen, "./whenthen.json")
}
