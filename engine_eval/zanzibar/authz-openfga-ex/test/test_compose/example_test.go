package test_compose

import (
	"context"
	"encoding/json"
	openfga "github.com/openfga/go-sdk"
	"github.com/stretchr/testify/assert"
	tc "github.com/testcontainers/testcontainers-go"
	"os"
	"testing"
)

func TestSomething(t *testing.T) {
	compose, err := tc.NewDockerCompose("docker-compose.yml")
	assert.NoError(t, err, "NewDockerComposeAPI()")

	t.Cleanup(func() {
		assert.NoError(t, compose.Down(context.Background(), tc.RemoveOrphans(true), tc.RemoveImagesLocal), "compose.Down()")
	})

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	assert.NoError(t, compose.Up(ctx, tc.Wait(true)), "compose.Up()")

	// do some testing here

	configuration, err := openfga.NewConfiguration(openfga.Configuration{
		ApiScheme: "http",
		ApiHost:   "0.0.0.0:8080",
	})

	if err != nil {
		t.Fatalf("%v", err)
	}

	apiClient := openfga.NewAPIClient(configuration)

	resp, _, err := apiClient.OpenFgaApi.CreateStore(context.Background()).Body(openfga.CreateStoreRequest{
		Name: openfga.PtrString("FGA Testcontainer Store"),
	}).Execute()

	if err != nil {
		t.Fatalf("Failed to create store. Error: %v", err)
	}
	t.Logf("Name: %s, ID: %s", resp.GetName(), resp.GetId())
	storeId := resp.GetId()

	content, err := os.ReadFile("model.json")

	if err != nil {
		t.Fatalf("failed to load model from file: %s", err.Error())
	}

	var body openfga.WriteAuthorizationModelRequest
	if err := json.Unmarshal(content, &body); err != nil {
		t.Errorf("Error unmarshalling: %v", err)
		// .. Handle error
		return
	}

	configuration, err = openfga.NewConfiguration(openfga.Configuration{
		ApiScheme: "http",
		ApiHost:   "0.0.0.0:8080",
		StoreId:   storeId,
	})

	if err != nil {
		t.Fatalf("%v", err)
	}

	apiClient = openfga.NewAPIClient(configuration)

	data, response, err := apiClient.OpenFgaApi.WriteAuthorizationModel(context.Background()).Body(body).Execute()
	if err != nil {
		t.Errorf("Error writing authorizationmodel: %v", err)
	}
	t.Logf("data: %v", data)
	t.Logf("response: %v", response)
}
