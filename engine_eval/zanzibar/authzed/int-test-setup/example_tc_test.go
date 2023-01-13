package int_test_setup

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	v1 "github.com/authzed/authzed-go/proto/authzed/api/v1"
	"github.com/authzed/authzed-go/v1"
	"github.com/authzed/grpcutil"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"testing"
)

type spicedbContainer struct {
	testcontainers.Container
	URI        string
	MappedPort string
}

func TestAuthzedTcExample(t *testing.T) {
	if testing.Short() {
		t.Skip("-test.short flag set, skipping integration test")
	}

	ctx := context.Background()
	spicedbContainer, err := setupSpiceDb(ctx, t)
	port := spicedbContainer.MappedPort

	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name   string
		schema string
	}{
		{
			"basic readback",
			`definition user {}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client, err := spicedbTestClient(port)

			if err != nil {
				t.Fatal(err)
			}

			_, err = client.WriteSchema(context.TODO(), &v1.WriteSchemaRequest{Schema: tt.schema})
			if err != nil {
				t.Fatal(err)
			}

			resp, err := client.ReadSchema(context.TODO(), &v1.ReadSchemaRequest{})
			if err != nil {
				t.Fatal(err)
			}

			if tt.schema != resp.SchemaText {
				t.Fatal(err)
			}
		})
	}
}

func setupSpiceDb(ctx context.Context, t *testing.T) (*spicedbContainer, error) {
	ctx = context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "authzed/spicedb:latest",
		ExposedPorts: []string{"50051/tcp", "50052/tcp"},
		WaitingFor:   wait.ForLog("grpc server started serving"),
		Cmd:          []string{"serve-testing"},
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
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
	defer func() {
		if err := container.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err.Error())
		}
	}()

	return &spicedbContainer{Container: container, URI: uri, MappedPort: mappedPort.Port()}, nil
}

// spicedbTestClient creates a new SpiceDB client with random credentials.
//
// The test server gives each set of a credentials its own isolated datastore
// so that tests can be run in parallel.
func spicedbTestClient(port string) (*authzed.Client, error) {
	// Generate a random credential to isolate this client from any others.
	buf := make([]byte, 20)
	if _, err := rand.Read(buf); err != nil {
		return nil, err
	}
	randomKey := base64.StdEncoding.EncodeToString(buf)

	return authzed.NewClient(
		"localhost:"+port,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpcutil.WithInsecureBearerToken(randomKey),
		grpc.WithBlock(),
	)
}
