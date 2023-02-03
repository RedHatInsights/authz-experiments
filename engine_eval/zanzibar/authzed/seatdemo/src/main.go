package main

import (
	"context"
	"log"

	v1 "github.com/authzed/authzed-go/proto/authzed/api/v1"
	"github.com/authzed/authzed-go/v1"
	"github.com/authzed/grpcutil"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	client, err := getSpiceDbApiClient("50051")

	if err != nil {
		log.Fatalf("unable to initialize client: %s", err)
	}

	checkConnection(client)

	if err != nil {
		log.Fatalf("unable to check for connection: %s", err)
	}

}

func getSpiceDbApiClient(port string) (*authzed.Client, error) {
	client, err := authzed.NewClient(
		"localhost:"+port,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpcutil.WithInsecureBearerToken("abcdefgh"),
	)
	return client, err
}

func checkConnection(client *authzed.Client) (schema string, err error) {
	ctx := context.Background()

	schemaResponse, err := client.ReadSchema(ctx, &v1.ReadSchemaRequest{})

	if err != nil {
		return
	}
	schema = schemaResponse.SchemaText
	log.Println(schema)
	return
}
