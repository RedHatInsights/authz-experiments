package handler

import (
	"context"
	"github.com/authzed/authzed-go/proto/authzed/api/v1"
	"github.com/labstack/echo/v4"
	"io"
	"log"
	"net/http"
)

type GetUserResponse struct {
	UserId string `json:"userid" xml:"userid"`
}

type GetUserResponses []GetUserResponse

func GetTenantUsers(c echo.Context) error {

	ctx := context.Background()
	if port == "" {
		port = "50051" //TODO
	}
	client, err := getSpiceDbApiClient(port)

	if err != nil {
		log.Fatalf("unable to initialize client: %s", err)
		return err
	}

	callingName := c.QueryParam("callingName")
	subject := &v1.SubjectReference{Object: &v1.ObjectReference{
		ObjectType: "user",
		ObjectId:   callingName,
	}}

	tenantId := c.Param("tenant")
	tenant := &v1.ObjectReference{
		ObjectType: "tenant",
		ObjectId:   tenantId,
	}

	log.Printf("Calling CheckPermission using Subject: %s for Resource: %s\n", subject.GetObject().ObjectId, tenant.GetObjectId())

	resp, err := client.CheckPermission(ctx, &v1.CheckPermissionRequest{
		Resource:   tenant,
		Permission: "is_owner",
		Subject:    subject,
	})

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "uh oh...")
	}

	if resp.Permissionship != v1.CheckPermissionResponse_PERMISSIONSHIP_HAS_PERMISSION {
		return echo.NewHTTPError(http.StatusForbidden, "nothing to see here (no really, there is, but you are not allowed to. ehehe)")
	}

	stream, err := client.LookupSubjects(ctx, &v1.LookupSubjectsRequest{
		Resource:          tenant,
		Permission:        "membership",
		SubjectObjectType: "user",
	})

	var result GetUserResponses
	for { //most likely weird and ugly, but works.
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err == nil {
			//now append to my structs
			result = append(result, GetUserResponse{
				UserId: resp.GetSubjectObjectId(),
			})
		}

		if err != nil {
			log.Fatal("tilt", err)
		}
	}
	return c.JSON(http.StatusOK, result)
}
