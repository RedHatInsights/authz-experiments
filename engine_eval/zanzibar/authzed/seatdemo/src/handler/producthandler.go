package handler

import (
	"context"
	"errors"
	v1 "github.com/authzed/authzed-go/proto/authzed/api/v1"
	"github.com/labstack/echo/v4"
	"io"
	"log"
	"net/http"
)

type ProductLicense struct {
	Name   string `json:"name" xml:"name"`
	Active int    `json:"active_licenses" xml:"active_licenses"`
	Max    int    `json:"max_seats" xml:"max_seats"`
}

func GetLicenseInfoForProductInstance(c echo.Context) error {
	licenseMap := map[string]int{"p1": 5, "p2": 4}

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

	resp, err := client.CheckPermission(ctx, &v1.CheckPermissionRequest{
		Resource:   tenant,
		Permission: "manage_seats",
		Subject:    subject,
	})

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "uh oh...")
	}

	//check permissions: is p1 tied to customer1? actual intent: is user allowed to access p1 licensing data? let's check manage_seats permission on tenant resource for now.
	if resp.Permissionship != v1.CheckPermissionResponse_PERMISSIONSHIP_HAS_PERMISSION {
		return echo.NewHTTPError(http.StatusForbidden, "You are not allowed to see licensing information. manage_seats is required (and too coarse grained, but for the sake of example it suffices.")
	}

	pInstance := c.Param("pinstance")
	productInstance := &v1.ObjectReference{
		ObjectType: "product_instance",
		ObjectId:   pInstance,
	}

	//if yes: get is_active_user LookupSubject for product_instance, count and put into active.
	stream, err := client.LookupSubjects(ctx, &v1.LookupSubjectsRequest{
		Resource:          productInstance,
		Permission:        "is_active_user",
		SubjectObjectType: "user",
	})

	var currentLicenseCount int
	for { //most likely weird and ugly, but works.
		_, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err == nil {
			//just count.
			currentLicenseCount++
		}

		if err != nil {
			log.Fatal("tilt", err)
		}
	}

	result := ProductLicense{
		Name:   pInstance,
		Active: currentLicenseCount,
		Max:    licenseMap[pInstance],
	}
	return c.JSON(http.StatusOK, result)
}

func GrantLicenseIfNotFull(c echo.Context) error {
	return errors.New("to be implemented")
}
