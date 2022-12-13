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
		StoreId:   "01GM3S2VB1H25YA3KKF3XNDJQ8",
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
	})
}
