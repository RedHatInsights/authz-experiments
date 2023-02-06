package handler

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/kinbiko/jsonassert"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestMain(m *testing.M) { //Acts as a test fixture for the package. NOTE: must be in a _test.go file.
	//Setup code runs before all tests
	db, err := getSpiceDbContainer()
	if err != nil {
		log.Fatalf("container not setup correctly: %s", err)
	}

	SetPort(db.MappedPort)

	m.Run() //All tests run here

	//Cleanup code will run after all tests in package
}

func get(uri string) *http.Request {
	return httptest.NewRequest(http.MethodGet, uri, strings.NewReader(""))
}

func post(uri string, body string) *http.Request {
	return reqWithBody(http.MethodPost, uri, body)
}

func put(uri string, body string) *http.Request {
	return reqWithBody(http.MethodPut, uri, body)
}

func delete(uri string) *http.Request {
	return httptest.NewRequest(http.MethodDelete, uri, strings.NewReader(""))
}

func reqWithBody(method string, uri string, body string) *http.Request {
	req := httptest.NewRequest(method, uri, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req
}

func runRequest(req *http.Request) *http.Response {
	e := GetEcho()
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	return rec.Result()
}

func assertHttpErrCodeAndMsg(t *testing.T, statusCode int, message string, resp *http.Response) {
	if assert.NotNil(t, resp) {
		assert.Equal(t, statusCode, resp.StatusCode)

		payload := new(strings.Builder)
		_, err := io.Copy(payload, resp.Body)
		assert.NoError(t, err)

		assert.Contains(t, payload.String(), message)
	}
}

func assertJsonResponse(t *testing.T, resp *http.Response, statusCode int, template string, args ...interface{}) {
	if assert.NotNil(t, resp) {
		assert.Equal(t, statusCode, resp.StatusCode)

		payload := new(strings.Builder)
		_, err := io.Copy(payload, resp.Body)
		assert.NoError(t, err)

		ja := jsonassert.New(t)
		ja.Assertf(payload.String(), template, args...)
	}
}

type spicedbContainer struct { //TODO: move out, instantiate only once etc
	testcontainers.Container
	URI        string
	MappedPort string
}

var spiceDB *spicedbContainer

func getSpiceDbContainer() (*spicedbContainer, error) {
	if spiceDB != nil {
		return spiceDB, nil
	}

	spiceDB, err := setupSpiceDb(context.Background())
	if err != nil {
		return nil, err
	}

	return spiceDB, nil
}

func setupSpiceDb(ctx context.Context) (*spicedbContainer, error) {
	ctx = context.Background()

	var (
		_, b, _, _ = runtime.Caller(0)
		basepath   = filepath.Dir(b)
	)

	req := testcontainers.ContainerRequest{
		Name:         "testspice",
		Image:        "authzed/spicedb:latest",
		ExposedPorts: []string{"50051/tcp", "50052/tcp"},
		WaitingFor:   wait.ForLog("grpc server started serving"),
		Mounts: testcontainers.Mounts(
			testcontainers.ContainerMount{
				Source: testcontainers.GenericBindMountSource{HostPath: path.Join(basepath, "../testresources/model.yaml")},
				Target: "/var/lib/spicedb/data/model.yaml"}),
		Env: map[string]string{
			"SPICEDB_GRPC_PRESHARED_KEY": "abcdefgh",
		},
		Cmd: []string{"serve", "--datastore-bootstrap-files", "/var/lib/spicedb/data/model.yaml"},
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
		Reuse:            true,
	})

	if err != nil {
		return nil, err
	}

	ip, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}

	mappedPort, err := container.MappedPort(ctx, "50051")

	if err != nil {
		return nil, err
	}

	uri := fmt.Sprintf("http://%s:%s", ip, mappedPort.Port())

	return &spicedbContainer{Container: container, URI: uri, MappedPort: mappedPort.Port()}, nil
}
