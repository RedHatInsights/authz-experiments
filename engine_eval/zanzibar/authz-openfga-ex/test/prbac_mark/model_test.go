package prbac_mark

import (
	"context"
	"encoding/json"
	"fmt"
	openfga "github.com/openfga/go-sdk"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"testing"
)

type TestDefinition struct {
	Name           string
	JsonResponse   string
	ResponseStatus int
	Method         string
	RequestPath    string
}

func TestWithOpenFGA(t *testing.T) {
	configuration, err := openfga.NewConfiguration(openfga.Configuration{
		ApiScheme: "http",
		ApiHost:   "0.0.0.0:8080",
		StoreId:   "01GM5CBY4TR8QJ97GNDJVF71TW",
	})
	if err != nil {
		t.Fatalf("%v", err)
	}

	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "openfga/openfga",
		ExposedPorts: []string{"8080/tcp", "3000/tcp"},
		WaitingFor:   wait.ForLog("HTTP server listening on '0.0.0.0:8080'..."),
		Cmd:          []string{"run"},
		Env: map[string]string{
			"OPENFGA_API_SCHEME": "http",
			"OPENFGA_API_HOST":   "0.0.0.0:8080",
			"OPENFGA_AUTH_MODEL": "{\"type_definitions\":[{\"type\":\"document\",\"relations\":{\"reader\":{\"this\":{}},\"writer\":{\"this\":{}},\"owner\":{\"this\":{}}}}]}",
		},
	}
	openFgaContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		t.Error(err)
	}
	defer func() {
		if err := openFgaContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err.Error())
		}
	}()

	apiClient := openfga.NewAPIClient(configuration)
	t.Run("Check", func(t *testing.T) {
		test := TestDefinition{
			Name:           "Check",
			JsonResponse:   `{"allowed":true, "resolution":""}`,
			ResponseStatus: 200,
			Method:         "POST",
			RequestPath:    "check",
		}
		requestBody := openfga.CheckRequest{
			TupleKey: &openfga.TupleKey{
				User:     openfga.PtrString("user:81684243-9356-4421-8fbf-a4f8d36aa31b"),
				Relation: openfga.PtrString("reader"),
				Object:   openfga.PtrString("document:roadmap"),
			},
		}

		var expectedResponse openfga.CheckResponse
		if err := json.Unmarshal([]byte(test.JsonResponse), &expectedResponse); err != nil {
			t.Fatalf("%v", err)
		}
		got, response, err := apiClient.OpenFgaApi.Check(context.Background()).Body(requestBody).Execute()
		if err != nil {
			t.Fatalf("%v", err)
		}
		fmt.Println(got)
		fmt.Println(response)

		if response.StatusCode != test.ResponseStatus {
			t.Fatalf("OpenFga%v().Execute() = %v, want %v", test.Name, response.StatusCode, test.ResponseStatus)
		}

		responseJson, err := got.MarshalJSON()
		if err != nil {
			t.Fatalf("%v", err)
		}

		if *got.Allowed != *expectedResponse.Allowed {
			t.Fatalf("OpenFga%v().Execute() = %v, want %v", test.Name, string(responseJson), test.JsonResponse)
		}

	})
}
