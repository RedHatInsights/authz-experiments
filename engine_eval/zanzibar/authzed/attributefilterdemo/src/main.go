package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	v1 "github.com/authzed/authzed-go/proto/authzed/api/v1"
	"github.com/authzed/authzed-go/v1"
	"github.com/authzed/grpcutil"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/structpb"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Provide a username and an app as commandline arguments. Ex: ./main alec advisor")
		return
	}

	userId := os.Args[1]
	app := os.Args[2]

	client, err := authzed.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpcutil.WithInsecureBearerToken("somerandomkeyhere"),
	)

	if err != nil {
		fmt.Printf("Unable to initialize client: %s\n", err)
		return
	}

	perms, err := getPermissions(client, userId, app)
	if err != nil {
		fmt.Printf("Failed to get permissions: %s\n", err)
		return
	}

	filters, err := getFilters(client, userId, app)
	if err != nil {
		fmt.Printf("Failed to get filters: %s\n", err)
		return
	}

	for _, perm := range perms {
		fmt.Printf("Permission: %s\n", decodePermission(perm))

		if filters, ok := filters[perm]; ok {
			fmt.Println("\tfilters:")

			for _, filterName := range filters {
				if filter, err := parseFilter(filterName); err != nil {
					fmt.Printf("error parsing filter: %s\n", err)
				} else {
					fmt.Printf("\t\t- %s\n", filter)
				}
			}
		}
	}
}

func getPermissions(client *authzed.Client, userId string, app string) ([]string, error) {
	ctx := context.Background()

	var objectIds []string

	lookupClient, err := client.LookupResources(ctx, &v1.LookupResourcesRequest{
		Subject:            userSubject(userId),
		Permission:         "granted",
		ResourceObjectType: "access",
		Context: &structpb.Struct{
			Fields: map[string]*structpb.Value{
				"current_app": structpb.NewStringValue(app),
			},
		},
	})

	if err != nil {
		return nil, err
	}

	for {
		response, err := lookupClient.Recv()

		switch {
		case errors.Is(err, io.EOF):
			return objectIds, nil
		case err != nil:
			return nil, err
		default:
			objectIds = append(objectIds, response.ResourceObjectId)
		}
	}
}

func getFilters(client *authzed.Client, userId string, app string) (map[string][]string, error) {
	ctx := context.Background()

	permissionsToFilters := make(map[string][]string)

	lookupClient, err := client.LookupResources(ctx, &v1.LookupResourcesRequest{
		Subject:            userSubject(userId),
		Permission:         "applies",
		ResourceObjectType: "filter",
		Context: &structpb.Struct{
			Fields: map[string]*structpb.Value{
				"current_app": structpb.NewStringValue(app),
			},
		},
	})

	if err != nil {
		return nil, err
	}

	for {
		response, err := lookupClient.Recv()

		switch {
		case errors.Is(err, io.EOF):
			return permissionsToFilters, nil
		case err != nil:
			return nil, err
		default:
			if permission, filter, err := parseFilterId(response.ResourceObjectId); err != nil {
				return nil, err
			} else {
				permissionsToFilters[permission] = append(permissionsToFilters[permission], filter)
			}
		}
	}
}

func parseFilterId(filterId string) (permission string, filter string, err error) {
	segments := strings.Split(filterId, "__")

	if len(segments) != 2 {
		return "", "", errors.New("malformed filterId %s - expected format permission_id_string__filter_id_string")
	}

	permission = segments[0]
	filter = segments[1]
	err = nil
	return
}

func decodePermission(permissionId string) string {
	str := strings.Replace(permissionId, "_", ":", -1)
	str = strings.Replace(str, "any", "*", -1)

	return str
}

func userSubject(id string) *v1.SubjectReference {
	return &v1.SubjectReference{
		Object: &v1.ObjectReference{
			ObjectType: "principal",
			ObjectId:   id,
		},
	}
}

type Filter struct {
	Key       string
	Operation string
	Value     string
}

func parseFilter(filterName string) (Filter, error) {
	segments := strings.Split(filterName, "_")

	if len(segments) != 3 {
		return Filter{}, errors.New("malformed filter name %s - expected format key_operation_value")
	}

	return Filter{
		Key:       segments[0],
		Operation: segments[1],
		Value:     segments[2],
	}, nil
}
