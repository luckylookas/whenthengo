package cleaning_pipe

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"strings"
	"testing"
)

func getTestString(count int) string {
	sb := strings.Builder{}
	sb.WriteString("[")
	for i := 0; i < count; i++ {
		sb.WriteString(
			`{
			"Data": "Content"
		},`)
	}
	test := strings.TrimSuffix(sb.String(), ",") + "]"
	return test
}

func getExpectedString(count int) string {
	sb := strings.Builder{}
	sb.WriteString("[")
	for i := 0; i < count; i++ {
		sb.WriteString(`{"data":"content"},`)
	}
	test := strings.TrimSuffix(sb.String(), ",") + "]"
	return test
}

var EMPTY = []byte{}
var tab = []byte("\t")
var cr = []byte("\r")
var nl = []byte("\n")
var space = []byte(" ")

func demoCleaner(stripped []byte) []byte {
	stripped = bytes.ReplaceAll(stripped, tab, EMPTY)
	stripped = bytes.ReplaceAll(stripped, cr, EMPTY)
	stripped = bytes.ReplaceAll(stripped, nl, EMPTY)
	stripped = bytes.ReplaceAll(stripped, space, EMPTY)
	return bytes.ToLower(stripped)
}

func TestCleanBodyPipe_Read(t *testing.T) {
	test := getTestString(10)
	expected := getExpectedString(10)

	stream := strings.NewReader(test)
	actual, err := ioutil.ReadAll(NewCleaningPipe(demoCleaner, stream))
	assert.NoError(t, err)
	assert.Equal(t, expected, string(actual))
}

func TestCleanBodyPipe_Read_EdgeCases(t *testing.T) {
	p := NewCleaningPipe(demoCleaner, strings.NewReader(""))
	actual, err := ioutil.ReadAll(p)
	assert.NoError(t, err)
	assert.ElementsMatch(t, actual, EMPTY, "empty reader")

	p = NewCleaningPipe(demoCleaner, nil)
	actual, err = ioutil.ReadAll(p)
	assert.NoError(t, err)
	assert.ElementsMatch(t, actual, EMPTY, "nil reader")

	p = NewCleaningPipe(demoCleaner, strings.NewReader(" "))
	actual, err = ioutil.ReadAll(p)
	assert.NoError(t, err)
	assert.ElementsMatch(t, actual, EMPTY, "just 1 whitespace")

	justemptys := strings.Repeat(" ", 10000) + "a"
	p = NewCleaningPipe(demoCleaner, strings.NewReader(justemptys))
	actual, err = ioutil.ReadAll(p)
	assert.NoError(t, err)
	assert.ElementsMatch(t, actual, []byte("a"), "full empty buffers before content")

	p = NewCleaningPipe(demoCleaner, strings.NewReader("    a"))
	actual, err = ioutil.ReadAll(p)
	assert.NoError(t, err)
	assert.Equal(t, string(actual), "a", "leading whitespace")

	p = NewCleaningPipe(demoCleaner, strings.NewReader("a  "))
	actual, err = ioutil.ReadAll(p)
	assert.NoError(t, err)
	assert.Equal(t, string(actual), "a", "trailing whitespace")

	p = NewCleaningPipe(demoCleaner, strings.NewReader("a    a"))
	actual, err = ioutil.ReadAll(p)
	assert.NoError(t, err)
	assert.Equal(t, string(actual), "aa", "prefix and suffix")
}