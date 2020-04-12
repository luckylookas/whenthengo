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
	"testing"
	"time"
)

func TestLatestReleaseWithTestContainersGo_Json(t *testing.T) {
	version := os.Getenv("RELEASE_VERSION")
	if version == "" {
		t.Log("no version set, cannot run integration tests")
		//t.FailNow()
	}
	wd, err := filepath.Abs("./")
	assert.NoError(t, err)

	volumemount := fmt.Sprintf("%s%ctest_resources", wd, os.PathSeparator)
	containerpath := fmt.Sprintf("%cin", os.PathSeparator)

	containerWhenThenFile := fmt.Sprintf("%s%c%s.json", containerpath, os.PathSeparator, t.Name())

	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "luckylukas/whenthengo:0.0.2",
			ExposedPorts: []string{"8080/tcp"},
			WaitingFor:    &wait.HTTPStrategy{
				Port:              "8080/tcp",
				Path:              "/whenthengoup",
				StatusCodeMatcher: func (status int) bool {
					return status == http.StatusOK
				}		,
				UseTLS:            false,
			},
			Env: map[string]string{
				"PORT": "8080",
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
	assert.NoError(t, err)

	err = container.Start(context.Background())


	assert.NoError(t, err)
	assert.NotNil(t, container.GetContainerID())

	//defer container.Terminate(ctx)
	ip, err := container.Host(ctx)
	assert.NoError(t, err)
	assert.NotEmpty(t, ip)


	port, err := container.MappedPort(ctx, "8080")
	assert.NoError(t, err)
	assert.NotEmpty(t, port.Port())


	httprequest, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s:%s/path/test", ip, port.Port()), nil)
	assert.NoError(t, err)
	time.Sleep(5 * time.Minute)
	httprequest.Header.Set("accept", "application/json")
	httprequest.Header.Set("unused", "ignored")

	start := time.Now()

	resp, err := http.DefaultClient.Do(httprequest)
	end := time.Now()
	assert.NoError(t, err)

	assert.True(t, (end.Nanosecond()/1000000 - start.Nanosecond()/1000000) > 2000)
	assert.Equal(t, 200, resp.StatusCode)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, "some-data", resp.Header.Get("some-header"))
	assert.Equal(t, "k", string(body))
}
