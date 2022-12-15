package openfgatc

import (
	"context"
	"encoding/json"
	openfga "github.com/openfga/go-sdk"
	"github.com/stretchr/testify/assert"
	tc "github.com/testcontainers/testcontainers-go"
	"net/http"
	"os"
	"testing"
)

func TestOpenFgaTcExample(t *testing.T) {
	setupComposeContainers(t)

	storeId := createOpenFgaStoreAndGetId(t)

	apiClient := getApiClientForStore(t, storeId)

	modelData := createAuthorizationModelFromJson(t, apiClient)

	createTuplesFromJson(t, apiClient) // returns the data but not needed here.

	createAssertionsFromJson(t, apiClient, modelData) // returns the data but not needed here.
	//TODO: actually check/assert sth with the given model
}

func createAssertionsFromJson(t *testing.T, apiClient *openfga.APIClient, modelData openfga.WriteAuthorizationModelResponse) *http.Response {
	assertions := readAssertionsFromFile(t)

	var assertionsReqBody openfga.WriteAssertionsRequest

	if err := json.Unmarshal(assertions, &assertionsReqBody.Assertions); err != nil {
		t.Errorf("Error unmarshalling assertions: %v", err)
	}

	assertionResponse, err := apiClient.OpenFgaApi.WriteAssertions(context.Background(), *modelData.AuthorizationModelId).Body(assertionsReqBody).Execute()

	if err != nil {
		t.Errorf("Error writing assertions: %v", err)
	}

	t.Logf("assertion resp: %v", assertionResponse)

	return assertionResponse
}

func createTuplesFromJson(t *testing.T, apiClient *openfga.APIClient) (map[string]interface{}, *http.Response) {
	tuples := readTuplesFromFile(t)
	var tuplesReqBody openfga.WriteRequest
	var tupleKeys openfga.TupleKeys
	if err := json.Unmarshal(tuples, &tupleKeys.TupleKeys); err != nil {
		t.Errorf("Error unmarshalling tuples: %v", err)
	}
	tuplesReqBody.SetWrites(tupleKeys)
	tupleData, tupleResponse, err := apiClient.OpenFgaApi.Write(context.Background()).Body(tuplesReqBody).Execute()

	if err != nil {
		t.Errorf("Error writing tuples: %v", err)
	}

	t.Logf("tuple resp data: %v", tupleData)
	t.Logf("tuple response: %v", tupleResponse)
	return tupleData, tupleResponse
}

func createAuthorizationModelFromJson(t *testing.T, apiClient *openfga.APIClient) openfga.WriteAuthorizationModelResponse {
	model := readModelFromFile(t)

	var modelReqBody openfga.WriteAuthorizationModelRequest

	if err := json.Unmarshal(model, &modelReqBody); err != nil {
		t.Errorf("Error unmarshalling model: %v", err)
	}

	modelData, modelResponse, err := apiClient.OpenFgaApi.WriteAuthorizationModel(context.Background()).Body(modelReqBody).Execute()

	if err != nil {
		t.Errorf("Error writing authorizationmodel: %v", err)
	}

	t.Logf("model resp data: %v", modelData)
	t.Logf("model response: %v", modelResponse)

	return modelData
}

func getApiClientForStore(t *testing.T, storeId string) *openfga.APIClient {
	configuration, err := openfga.NewConfiguration(openfga.Configuration{
		ApiScheme: "http",
		ApiHost:   "0.0.0.0:8080",
		StoreId:   storeId,
	})

	if err != nil {
		t.Errorf("Error creating new configuration for store: %v Error: %v", storeId, err)
	}
	apiClient := openfga.NewAPIClient(configuration)
	return apiClient
}

func createOpenFgaStoreAndGetId(t *testing.T) string {
	configuration, err := openfga.NewConfiguration(openfga.Configuration{
		ApiScheme: "http",
		ApiHost:   "0.0.0.0:8080",
	})

	if err != nil {
		t.Fatalf("%v", err)
	}

	apiClient := openfga.NewAPIClient(configuration)

	resp, _, er := apiClient.OpenFgaApi.CreateStore(context.Background()).Body(openfga.CreateStoreRequest{
		Name: "OpenFGA Testcontainer Store",
	}).Execute()

	if er != nil {
		t.Fatalf("Failed to create store. Error: %v", err)
	}

	t.Logf("Name: %s, ID: %s", resp.GetName(), resp.GetId())
	storeId := resp.GetId()
	return storeId
}

func readModelFromFile(t *testing.T) []byte {
	model, err := os.ReadFile("model.json")

	if err != nil {
		t.Fatalf("failed to load model from file: %s", err.Error())
	}
	return model
}

func readTuplesFromFile(t *testing.T) []byte {
	tuples, err := os.ReadFile("tuples.json")

	if err != nil {
		t.Fatalf("failed to load tuples from file: %s", err.Error())
	}
	return tuples
}

func readAssertionsFromFile(t *testing.T) []byte {
	model, err := os.ReadFile("assertions.json")

	if err != nil {
		t.Fatalf("failed to load assertions from file: %s", err.Error())
	}
	return model
}

func setupComposeContainers(t *testing.T) {
	compose, err := tc.NewDockerCompose("docker-compose.yml")
	assert.NoError(t, err, "NewDockerComposeAPI()")

	t.Cleanup(func() {
		assert.NoError(t, compose.Down(context.Background(), tc.RemoveOrphans(true), tc.RemoveImagesLocal), "compose.Down()")
	})

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	assert.NoError(t, compose.Up(ctx, tc.Wait(true)), "compose.Up()")
}
