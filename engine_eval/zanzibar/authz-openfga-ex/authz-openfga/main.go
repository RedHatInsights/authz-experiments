package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/golang/glog"
	openfga "github.com/openfga/go-sdk"
	"os"
)

func main() {
	configuration, err := openfga.NewConfiguration(openfga.Configuration{
		//ApiScheme: os.Getenv("OPENFGA_API_SCHEME"), // optional, defaults to "https"
		//ApiHost:   os.Getenv("OPENFGA_API_HOST"),
		ApiScheme: "http",         // Optional. Can be "http" or "https". Defaults to "https"
		ApiHost:   "0.0.0.0:8080", // required, define without the scheme (e.g. api.openfga.example instead of https://api.openfga.example)
	})
	AUTH_MODEL := os.Getenv("OPENFGA_AUTH_MODEL")

	_ = flag.CommandLine.Parse([]string{})

	//pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	// Always log to stderr by default
	if err := flag.Set("logtostderr", "true"); err != nil {
		glog.Infof("Unable to set logtostderr to true")
	}
	if err != nil {
		glog.Errorf("Error:%s", err)
	}

	apiClient := openfga.NewAPIClient(configuration)
	resp, _, err := apiClient.OpenFgaApi.ListStores(context.Background()).Execute()
	if err != nil {
		glog.Errorf("Error:%s", err)
	}

	if len(*resp.Stores) > 0 {
		for _, store := range *resp.Stores {
			if store.GetName() == "authz_usermgmt" {
				glog.Info("Setting store Id")
				apiClient.SetStoreId(store.GetId())
				glog.Infof("Store name:%s", store.GetName())
				glog.Infof("Store id:%s", store.GetId())
			}
		}
	} else {
		resp, err := createStore(apiClient)
		resp.GetId()
		apiClient.SetStoreId(resp.GetId())
		if err != nil {
			glog.Errorf("Error:%s", err)
		}
		glog.Infof("Store created:%s", resp.GetName())
	}

	//var writeAuthorizationModelRequestString = "{\n  \"type_definitions\": [\n    {\n      \"type\": \"user\"\n    },\n    {\n      \"type\": \"group\",\n      \"relations\": {\n        \"member\": {\n          \"this\": {}\n        }\n      }\n    },\n    {\n      \"type\": \"resource\",\n      \"relations\": {\n        \"writer\": {\n          \"this\": {}\n        },\n        \"reader\": {\n          \"union\": {\n            \"child\": [\n              {\n                \"this\": {}\n              },\n              {\n                \"computedUserset\": {\n                  \"relation\": \"writer\"\n                }\n              }\n            ]\n          }\n        }\n      }\n    }\n  ],\n  \"schema_version\": \"1.0\"\n}"
	var writeAuthorizationModelRequestString = AUTH_MODEL
	var body openfga.WriteAuthorizationModelRequest
	if err := json.Unmarshal([]byte(writeAuthorizationModelRequestString), &body); err != nil {
		glog.Errorf("Error :%s", err)
		return
	}

	data, response, err := apiClient.OpenFgaApi.WriteAuthorizationModel(context.Background()).Body(body).Execute()
	if err != nil {
		glog.Errorf("Error :%s", err)
	}
	//
	fmt.Println(data.GetAuthorizationModelId())
	//
	fmt.Println(response.Status)

}

func createStore(apiClient *openfga.APIClient) (openfga.CreateStoreResponse, error) {
	resp, _, err := apiClient.OpenFgaApi.CreateStore(context.Background()).Body(openfga.CreateStoreRequest{
		Name: openfga.PtrString("authz_usermgmt"),
	}).Execute()
	return resp, err
}
