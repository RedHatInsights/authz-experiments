package testsa

import (
	"authz-openfga-ex/utils"
	"context"
	"encoding/json"
	"fmt"
	"testing"

	openfga "github.com/openfga/go-sdk"
)

type TestDefinition struct {
	Name           string
	JsonResponse   string
	ResponseStatus int
	Method         string
	RequestPath    string
}

func readStoreId() (string, error) {
	var storeid string
	var err error
	file := "../storeId.txt"
	err = utils.ReadFileValueString(file, &storeid)
	return storeid, err
}

func TestOpenFgaApiConfigurationSAEx(t *testing.T) {
	t.Run("Providing no store id should not error", func(t *testing.T) {
		_, err := openfga.NewConfiguration(openfga.Configuration{
			ApiHost: "api.fga.example",
		})

		if err != nil {
			t.Fatalf("%v", err)
		}
	})
}

func TestOpenFgaApiSAExReadAccess(t *testing.T) {
	storeid, err := readStoreId()
	configuration, err := openfga.NewConfiguration(openfga.Configuration{
		ApiScheme: "http",
		ApiHost:   "0.0.0.0:8080",
		StoreId:   storeid,
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
				User:     openfga.PtrString("srvcaccnt1-aspian"),
				Relation: openfga.PtrString("read"),
				Object:   openfga.PtrString("metrics:read"),
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

func TestOpenFgaApiSAExListSAPermissions(t *testing.T) {
	storeid, err := readStoreId()
	configuration, err := openfga.NewConfiguration(openfga.Configuration{
		ApiScheme: "http",
		ApiHost:   "0.0.0.0:8080",
		StoreId:   storeid,
	})

	if err != nil {
		t.Fatalf("%v", err)
	}

	apiClient := openfga.NewAPIClient(configuration)

	t.Run("List", func(t *testing.T) {
		test := TestDefinition{
			Name:           "list-objects",
			JsonResponse:   `{"continuation_token":"","tuples":[{"key":{"object":"metrics:read","relation":"read","user":"group:aspian-metrics-read#member"},"timestamp":"2022-12-15T18:45:38.22683Z"}]}`,
			ResponseStatus: 200,
			Method:         "POST",
			RequestPath:    "list-objects",
		}
		//authModelID := "01GMBZNK7VA1CBHPRQRCZGYYZB"
		requestBody := openfga.ReadRequest{
			//AuthorizationModelId: &authModelID,
			TupleKey: &openfga.TupleKey{
				User:     openfga.PtrString("group:aspian-metrics-read#member"),
				Relation: openfga.PtrString(""),
				Object:   openfga.PtrString("metrics:read"),
			},
		}

		var expectedResponse openfga.ReadResponse
		if err := json.Unmarshal([]byte(test.JsonResponse), &expectedResponse); err != nil {
			t.Fatalf("%v", err)
		}
		got, response, err := apiClient.OpenFgaApi.Read(context.Background()).Body(requestBody).Execute()
		if err != nil {
			t.Fatalf("%v", err)
		}
		// fmt.Println(got)
		// fmt.Println(response)

		fmt.Println(got.GetTuples())
		if response.StatusCode != test.ResponseStatus {
			t.Fatalf("OpenFga%v().Execute() = %v want %v", test.Name, response.StatusCode, test.ResponseStatus)
		}

		responseJson, err := got.MarshalJSON()
		if err != nil {
			t.Fatalf("%v", err)
		}

		if len(*got.Tuples) != len(*expectedResponse.Tuples) {
			t.Fatalf("OpenFga%v().Execute() = %v, want %v", test.Name, string(responseJson), test.JsonResponse)
		}

	})
}
