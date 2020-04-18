// build integration

package client

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

func TestNewClient_Integration(t *testing.T) {
	version := os.Getenv("RELEASE_VERSION")
	if version == "" {
		t.Log("no version set, cannot run integration tests")
		t.FailNow()
	}

	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "luckylukas/whenthengo:" + version,
			ExposedPorts: []string{"80/tcp"},
			WaitingFor:    wait.ForHTTP("/whenthengo/up"),
		},
		Started: true,
	}
	ctx := context.Background()
	container, err := testcontainers.GenericContainer(ctx, req)
	if err != nil {
		if strings.Contains(err.Error(),  "Cannot connect to the Docker daemon") {
			t.Log("docker socket error when starting container - this is a known issue and does not fail the test")
		} else {
			t.Fatal(err)
		}
	}

	defer container.Terminate(ctx)
	ip, err := container.Host(ctx)
	assert.NoError(t, err)
	assert.NotEmpty(t, ip)

	port, err := container.MappedPort(ctx, "80")
	assert.NoError(t, err)
	assert.NotEmpty(t, port.Port())

	err = NewClient(ip, port.Port()).
		WhenRequest().
		ThenReply().
		And().
		WhenRequest().
		WithUrl("/data/").
		WithMethod("post").
		ThenReply().
		WithDelay(2000).
		WithStatus(201).
		AndDo().
		Publish(ctx)

	assert.NoError(t, err)

	resp, err := http.Get("http://" + ip + ":" + port.Port() + "/")
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	start := time.Now()
	resp, err = http.Post("http://" + ip + ":" + port.Port() + "/data", "application/json", nil)
	assert.True(t, time.Since(start) > 1800 * time.Millisecond)
	assert.NoError(t, err)
	assert.Equal(t, 201, resp.StatusCode)
}
