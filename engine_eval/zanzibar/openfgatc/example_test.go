package openfgatc

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"

	openfga "github.com/openfga/go-sdk"
	"github.com/stretchr/testify/assert"
	tc "github.com/testcontainers/testcontainers-go"
)

func TestOpenFgaTcExample(t *testing.T) {
	setupComposeContainers(t)

	storeId := createOpenFgaStoreAndGetId(t)

	apiClient := getApiClientForStore(t, storeId)

	modelData := createAuthorizationModelFromJson(t, apiClient)

	createTuplesFromJson(t, apiClient) // returns the data but not needed here.

	assertions := createAssertionsFromJson(t, apiClient, modelData) // returns the data but not needed here.

	checkAssertions(t, apiClient, modelData.AuthorizationModelId, assertions) //Checks the assertions against the model
	//TODO: actually check/assert sth with the given model
}

func checkAssertions(t *testing.T, apiClient *openfga.APIClient, modelId *string, assertions []openfga.Assertion) {
	trace := false

	for _, assertion := range assertions {
		body := openfga.CheckRequest{TupleKey: assertion.TupleKey, ContextualTuples: nil, AuthorizationModelId: modelId, Trace: &trace}
		result, _, err := apiClient.OpenFgaApi.Check(context.Background()).Body(body).Execute()

		if err != nil {
			t.Errorf("Error checking assertion tuple (%s): %v", tupleKeyToString(assertion.TupleKey), err)
			return
		}

		if assertion.Expectation != *result.Allowed {
			t.Errorf("Assertion failed! %s - Expected: %t, Actual: %t", tupleKeyToString(assertion.TupleKey), assertion.Expectation, *result.Allowed)
		}
	}
}

func createAssertionsFromJson(t *testing.T, apiClient *openfga.APIClient, modelData openfga.WriteAuthorizationModelResponse) []openfga.Assertion {
	jsonData := readAssertionsFromFile(t)

	var assertions []openfga.Assertion

	if err := json.Unmarshal(jsonData, &assertions); err != nil {
		t.Errorf("Error unmarshalling assertions: %v", err)
	}

	assertionResponse, err := apiClient.OpenFgaApi.WriteAssertions(context.Background(), *modelData.AuthorizationModelId).Body(*openfga.NewWriteAssertionsRequest(assertions)).Execute()

	if err != nil {
		t.Errorf("Error writing assertions: %v", err)
	}

	t.Logf("assertion resp: %v", assertionResponse)

	return assertions
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

func tupleKeyToString(val openfga.TupleKey) string {
	return fmt.Sprintf("User: %s, Relation: %s, Object: %s", *val.User, *val.Relation, *val.Object)
}
