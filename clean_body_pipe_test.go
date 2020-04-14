package main

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"strings"
	"testing"
)

func TestCleanBodyPipe_Read_Long(t *testing.T) {
	sb := strings.Builder{}
	for i:=0; i < 100000; i++ {
		sb.WriteString("aA\n")
	}
	test := sb.String()

	sb = strings.Builder{}
	for i:=0; i < 100000; i++ {
		sb.WriteString("aa")
	}
	expected := sb.String()

	stream := strings.NewReader(test)
	actual, err := ioutil.ReadAll(CleanBodyPipe{stream})
	assert.NoError(t, err)
	assert.Equal(t, expected, string(actual))
}

func TestCleanBodyPipe_Read_Simple(t *testing.T) {
	stream := strings.NewReader("Aa\n a")
	actual, err := ioutil.ReadAll(CleanBodyPipe{stream})
	assert.NoError(t, err)
	assert.Equal(t, "aaa", string(actual))
}