package test

import (
	"context"
	"encoding/json"
	"fmt"
	openfga "github.com/openfga/go-sdk"
	"testing"
)

type TestDefinition struct {
	Name           string
	JsonResponse   string
	ResponseStatus int
	Method         string
	RequestPath    string
}

func TestOpenFgaApiConfiguration(t *testing.T) {
	t.Run("Providing no store id should not error", func(t *testing.T) {
		_, err := openfga.NewConfiguration(openfga.Configuration{
			ApiHost: "api.fga.example",
		})

		if err != nil {
			t.Fatalf("%v", err)
		}
	})
}

func TestOpenFgaApi(t *testing.T) {
	configuration, err := openfga.NewConfiguration(openfga.Configuration{
		ApiScheme: "http",
		ApiHost:   "0.0.0.0:8080",
		StoreId:   "01GM5CBY4TR8QJ97GNDJVF71TW",
	})
	if err != nil {
		t.Fatalf("%v", err)
	}

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
