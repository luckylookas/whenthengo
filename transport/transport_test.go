package transport

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewWhenThenGoHttpClient(t *testing.T) {
	config := NewClient("", "").WhenRequest().WithUrl("/path").WithMethod("get").ThenReply().WithStatus(201).AndDo().Return()
	client := NewWhenThenGoHttpClient(config)

	resp, err := client.Get("http://anyhost:90/path")

	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, 201)
}
