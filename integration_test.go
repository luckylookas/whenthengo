// +build integration

package main

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestLatestReleaseWithTestContainersGo_Json(t *testing.T) {
	version := os.Getenv("RELEASE_VERSION")
	if version == "" {
		t.Log("no version set, cannot run integration tests")
		t.FailNow()
	}
	wd, err := filepath.Abs("./")
	assert.NoError(t, err)

	volumemount := fmt.Sprintf("%s%ctest_resources", wd, os.PathSeparator)
	containerpath := fmt.Sprintf("%cin", os.PathSeparator)

	containerWhenThenFile := fmt.Sprintf("%s%c%s.json", containerpath, os.PathSeparator, t.Name())

	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "luckylukas/whenthengo:" + version,
			ExposedPorts: []string{"80/tcp"},
			WaitingFor:    &wait.HTTPStrategy{
				Port:              "80/tcp",
				Path:              "/whenthengoup",
				StatusCodeMatcher: func (status int) bool {
					return status == http.StatusOK
				}		,
				UseTLS:            false,
			},
			Env: map[string]string{
				"PORT": "80",
				"WHENTHEN": containerWhenThenFile,
			},
			BindMounts: map[string]string{
				volumemount: containerpath,
			},
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


	// GET
	httprequest, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s:%s/path/test", ip, port.Port()), nil)
	assert.NoError(t, err)
	httprequest.Header.Set("accept", "Application/json")
	httprequest.Header.Set("unused", "ignored")

	start := time.Now()

	resp, err := http.DefaultClient.Do(httprequest)
	assert.NoError(t, err)

	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	defer resp.Body.Close()

	assert.True(t, time.Since(start).Milliseconds() > 1900)
	assert.Equal(t, 200, resp.StatusCode)
	assert.NoError(t, err)
	assert.Equal(t, "some-data", resp.Header.Get("some-header"))
	assert.Equal(t, "k", string(body))


	// POST with whitepsace ignoring Body matcher
	httprequest, err = http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s:%s/path/test",
		ip, port.Port()),
		strings.NewReader(`{"data":"content"}`))
	assert.NoError(t, err)
	httprequest.Header.Set("accept", "Application/json")
	httprequest.Header.Set("unused", "ignored")

	start = time.Now()

	resp, err = http.DefaultClient.Do(httprequest)
	assert.NoError(t, err)

	defer resp.Body.Close()

	assert.Equal(t, 201, resp.StatusCode)
	assert.Equal(t, "some-data", resp.Header.Get("some-header"))

	// NO MATCH
	httprequest, err = http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s:%s/path/test",
		ip, port.Port()),
		strings.NewReader(`{"data":"other"}`))
	assert.NoError(t, err)

	start = time.Now()

	resp, err = http.DefaultClient.Do(httprequest)
	assert.NoError(t, err)

	defer resp.Body.Close()

	assert.Equal(t, 404, resp.StatusCode)

}
