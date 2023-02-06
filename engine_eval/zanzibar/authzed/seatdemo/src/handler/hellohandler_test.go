package handler

import (
	"context"
	"fmt"
	"net/http"
	"path"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestHelloHandler(t *testing.T) {
	ctx := context.Background()
	db, err := setupSpiceDb(ctx, t)
	if err != nil {
		t.Fatalf("container not setup correctly: %s", err)
	}

	SetPort(db.MappedPort)

	resp := runRequest(get("/"))
	assertHttpErrCodeAndMsg(t, http.StatusOK, "Connection to spiceDB successfully established", resp)
}

type spicedbContainer struct { //TODO: move out, instantiate only once etc
	testcontainers.Container
	URI        string
	MappedPort string
}

func setupSpiceDb(ctx context.Context, t *testing.T) (*spicedbContainer, error) {
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
