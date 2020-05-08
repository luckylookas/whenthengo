package transport

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewClient_Defaults(t *testing.T) {
	whenthen := NewClient("localhost", "8080").
		WhenRequest().
		ThenReply().
		And().
		WhenRequest().
		WithMethod("post").
		ThenReply().
		AndDo().
		Return()

	assert.Equal(t, whenthen[0].When.Url, "/")
	assert.Equal(t, whenthen[0].When.Method, "get")
	assert.Equal(t, whenthen[0].Then.Status, 200)
	assert.Equal(t, whenthen[1].When.Method, "post")

}
