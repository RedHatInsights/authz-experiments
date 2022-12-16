package main

import (
	"authz-openfga-ex/utils"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/golang/glog"
	openfga "github.com/openfga/go-sdk"
)

func main() {
	configuration, err := openfga.NewConfiguration(openfga.Configuration{
		//ApiScheme: os.Getenv("OPENFGA_API_SCHEME"), // optional, defaults to "https"
		//ApiHost:   os.Getenv("OPENFGA_API_HOST"),
		ApiScheme: "http",         // Optional. Can be "http" or "https". Defaults to "https"
		ApiHost:   "0.0.0.0:8080", // required, define without the scheme (e.g. api.openfga.example instead of https://api.openfga.example)
	})

	AUTH_MODEL_JSON_FILE := os.Getenv("OPENFGA_AUTH_MODEL_JSON_FILE")
	AUTH_MODEL_TUPLE_JSON_FILE := os.Getenv("OPENFGA_AUTH_MODEL_TUPLES_JSON_FILE")
	AUTH_MODEL_ASSERTION_JSON_FILE := os.Getenv("OPENFGA_AUTH_MODEL_ASSERTION_JSON_FILE")

	_ = flag.CommandLine.Parse([]string{})

	//pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	// Always log to stderr by default
	if err := flag.Set("logtostderr", "true"); err != nil {
		glog.Infof("Unable to set logtostderr to true")
	}
	if err != nil {
		glog.Errorf("Error:%s", err)
	}

	// Preflight checks
	if utils.IsNil(AUTH_MODEL_JSON_FILE) {
		glog.Fatalf("Auth Model JSON File not specified. Please set env value `OPENFGA_AUTH_MODEL_JSON_FILE`")
	}

	apiClient := openfga.NewAPIClient(configuration)
	resp, _, err := apiClient.OpenFgaApi.ListStores(context.Background()).Execute()
	if err != nil {
		glog.Errorf("Error:%s", err)
	}

	var storeId string
	if len(*resp.Stores) > 0 {
		for _, store := range *resp.Stores {
			if store.GetName() == "authz_usermgmt" {
				glog.Info("Setting store Id")
				storeId = store.GetId()

				glog.Infof("Store name:%s", store.GetName())
				glog.Infof("Store id:%s", store.GetId())
			}
		}
	} else {
		resp, err := createStore(apiClient)
		storeId = resp.GetId()
		if err != nil {
			glog.Errorf("Error:%s", err)
		}
		glog.Infof("Store created:%s", resp.GetName())
	}

	//Set the storeID in the client
	if utils.IsNotNil(storeId) {
		apiClient.SetStoreId(storeId)
		utils.CreateFileFromStringData("storeId.txt", storeId) // Used in testing to get the storeID
	}

	authzModelID, err := createAuthModel(apiClient, AUTH_MODEL_JSON_FILE)
	if err != nil {
		glog.Fatalf("Error creating Authmodel: %s", err)
	}

	err = createTuples(apiClient, AUTH_MODEL_TUPLE_JSON_FILE)
	if err != nil {
		glog.Errorf("Error :%s", err)
	}

	err = createAssertions(apiClient, authzModelID, AUTH_MODEL_ASSERTION_JSON_FILE)
	if err != nil {
		glog.Errorf("Error :%s", err)
	}

}

func createStore(apiClient *openfga.APIClient) (openfga.CreateStoreResponse, error) {
	resp, _, err := apiClient.OpenFgaApi.CreateStore(context.Background()).Body(openfga.CreateStoreRequest{
		Name: openfga.PtrString("authz_usermgmt"),
	}).Execute()
	return resp, err
}

func createTuples(apiClient *openfga.APIClient, tuplesJsonFilePath string) error {

	jsonAuthModelTuplesFile, err := os.Open(tuplesJsonFilePath)
	if err != nil {
		glog.Errorf("Error :%s", err)
		return err
	}
	defer jsonAuthModelTuplesFile.Close()
	byteValue, err := ioutil.ReadAll(jsonAuthModelTuplesFile)
	if err != nil {
		glog.Errorf("Error :%s", err)
		return err
	}

	var body openfga.WriteRequest

	if err := json.Unmarshal(byteValue, &body); err != nil {
		glog.Errorf("Error :%s", err)
		return err
	}

	_, response, err := apiClient.OpenFgaApi.Write(context.Background()).Body(body).Execute()
	if err != nil {
		glog.Errorf("Error :%s", err)
	}
	fmt.Println(response.Status)
	return err
}

func createAssertions(apiClient *openfga.APIClient, authorizationModelId, assertionsJsonFilePath string) error {
	jsonAuthModelAssertionFile, err := os.Open(assertionsJsonFilePath)
	if err != nil {
		glog.Errorf("Error :%s", err)
		return err
	}
	defer jsonAuthModelAssertionFile.Close()
	byteValue, err := ioutil.ReadAll(jsonAuthModelAssertionFile)
	if err != nil {
		glog.Errorf("Error :%s", err)
		return err
	}

	var body openfga.WriteAssertionsRequest

	if err := json.Unmarshal(byteValue, &body); err != nil {
		glog.Errorf("Error :%s", err)
		return err
	}

	response, err := apiClient.OpenFgaApi.WriteAssertions(context.Background(), authorizationModelId).Body(body).Execute()
	if err != nil {
		glog.Errorf("Error :%s", err)
	}
	fmt.Println(response.Status)
	return err

}

func createAuthModel(apiClient *openfga.APIClient, authmodelJsonFilePath string) (string, error) {

	jsonAuthModelFile, err := os.Open(authmodelJsonFilePath)
	if err != nil {
		glog.Errorf("Error :%s", err)
		return "", err
	}
	defer jsonAuthModelFile.Close()
	byteValue, err := ioutil.ReadAll(jsonAuthModelFile)
	if err != nil {
		glog.Errorf("Error :%s", err)
		return "", err
	}

	var body openfga.WriteAuthorizationModelRequest

	if err := json.Unmarshal(byteValue, &body); err != nil {
		glog.Errorf("Error :%s", err)
		return "", err
	}

	data, response, err := apiClient.OpenFgaApi.WriteAuthorizationModel(context.Background()).Body(body).Execute()
	if err != nil {
		glog.Errorf("Error :%s", err)
		return "", err
	}

	utils.CreateFileFromStringData("authmodel.txt", data.GetAuthorizationModelId())
	glog.Infof("Create Authmodel response code: %s", response.Status)
	return data.GetAuthorizationModelId(), err

}

func GetListOfObjects(apiClient *openfga.APIClient, user, relation, object string) {
	requestBody := openfga.ReadRequest{
		TupleKey: &openfga.TupleKey{
			User:     openfga.PtrString(user),
			Relation: openfga.PtrString(relation),
			Object:   openfga.PtrString(object),
		},
	}

	data, response, err := apiClient.OpenFgaApi.Read(context.Background()).Body(requestBody).Execute()
	if err != nil {
		glog.Errorf("Error :%s", err)
	}
	fmt.Println(response.Status)
	tuples := data.GetTuples()
	fmt.Println(tuples)
}
