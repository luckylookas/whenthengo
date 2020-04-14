// +build integration

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
		WithUri("/data/").
		WithMethod("post").
		ThenReply().
		WithStatus(302).
		AndDo().
		Publish(ctx)

	assert.NoError(t, err)

	resp, err := http.Get("http://" + ip + ":" + port.Port() + "/")
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	resp, err = http.Post("http://" + ip + ":" + port.Port() + "/data", "application/json", nil)
	assert.NoError(t, err)
	assert.Equal(t, 302, resp.StatusCode)
}
