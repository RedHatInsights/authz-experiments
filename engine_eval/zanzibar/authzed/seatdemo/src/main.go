package main

import (
	"context"
	"log"
	"net/http"

	v1 "github.com/authzed/authzed-go/proto/authzed/api/v1"
	"github.com/authzed/authzed-go/v1"
	"github.com/authzed/grpcutil"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CheckConnectionResponse struct {
	Message string `json:"message" xml:"message"`
	Schema  string `json:"schema" xml:"schema"`
}

var port string

func main() {
	// Echo instance
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", hello)

	// Start server
	e.Logger.Fatal(e.Start(":3000"))
}

func hello(c echo.Context) error {
	client, err := getSpiceDbApiClient(port)

	if err != nil {
		log.Fatalf("unable to initialize client: %s", err)
	}

	schema, err := checkSpiceDbConnection(client)

	if err != nil {
		c.Error(err)
		return err
	}

	ccresp := &CheckConnectionResponse{
		Message: "Connection to spiceDB successfully established!",
		Schema:  schema,
	}
	return c.JSON(http.StatusOK, ccresp)

}

func getSpiceDbApiClient(port string) (*authzed.Client, error) {
	client, err := authzed.NewClient(
		"localhost:"+port,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpcutil.WithInsecureBearerToken("abcdefgh"),
	)
	return client, err
}

func checkSpiceDbConnection(client *authzed.Client) (schema string, err error) {
	ctx := context.Background()

	schemaResponse, err := client.ReadSchema(ctx, &v1.ReadSchemaRequest{})

	if err != nil {
		return
	}
	schema = schemaResponse.SchemaText
	log.Println(schema)
	return
}

func setPort(p string) {
	port = p
}
