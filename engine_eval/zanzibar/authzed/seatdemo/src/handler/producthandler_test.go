package handler

import (
	"context"
	"github.com/kinbiko/jsonassert"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGrantLicenseReturns403ForUserNotMemberOfTenant(t *testing.T) {
	ctx := context.Background()
	/*using one tc instance per test bc don't quite know how to create fixtures and stuff. concious technical debt for now. enlighten me ;)*/
	db, err := setupSpiceDb(ctx, t)
	if err != nil {
		t.Fatalf("container not setup correctly: %s", err)
	}

	SetPort(db.MappedPort)
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/tenant/customer1/product/:pinstance/license", strings.NewReader("userId=t2user3"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded") //Q: are there no constants for this in go? //TODO: evaluate why binding in code does not return error when this header is not set...
	rec := httptest.NewRecorder()

	echoCtx := e.NewContext(req, rec)
	echoCtx.SetPath("/tenant/:tenant/product/:pinstance/license")
	echoCtx.SetParamNames("tenant", "pinstance")
	echoCtx.SetParamValues("customer1", "p1")

	err2 := GrantLicenseIfNotFull(echoCtx)
	if assert.NotNil(t, err2) {
		he, ok := err2.(*echo.HTTPError)
		if ok {
			assert.Equal(t, http.StatusForbidden, he.Code)
			assert.Contains(t, he.Message, "User t2user3 is not a member of licensed tenant customer1")
		}
	}
}

func TestGrantLicenseRevokesGrantIfMaxReached(t *testing.T) {
	ctx := context.Background()
	/*using one tc instance per test bc don't quite know how to create fixtures and stuff. concious technical debt for now. enlighten me ;)*/
	db, err := setupSpiceDb(ctx, t)
	if err != nil {
		t.Fatalf("container not setup correctly: %s", err)
	}

	SetPort(db.MappedPort)
	e := echo.New()

	req := httptest.NewRequest(http.MethodPost, "/tenant/customer1/product/:pinstance/license", strings.NewReader("userId=t2user3"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded") //Q: are there no constants for this in go? //TODO: evaluate why binding in code does not return error when this header is not set...
	rec := httptest.NewRecorder()

	echoCtx := e.NewContext(req, rec)
	echoCtx.SetPath("/tenant/:tenant/product/:pinstance/license")
	echoCtx.SetParamNames("tenant", "pinstance")
	echoCtx.SetParamValues("customer2", "p2")

	//weird weird way of testing for errors in echo
	if assert.NoError(t, GrantLicenseIfNotFull(echoCtx)) {
		assert.Equal(t, http.StatusConflict, rec.Code)
		ja := jsonassert.New(t)
		// find some sort of payload
		granted := false
		reason := "Maximum seats exceeded. Please extend your license."
		ja.Assertf(rec.Body.String(), `
	{
		"granted": %t,
		"reason": "%s"
	}`, granted, reason)
	}
}

func TestGrantLicenseReturnsBadRequestWithoutBody(t *testing.T) {
	ctx := context.Background()
	/*using one tc instance per test bc don't quite know how to create fixtures and stuff. concious technical debt for now. enlighten me ;)*/
	db, err := setupSpiceDb(ctx, t)
	if err != nil {
		t.Fatalf("container not setup correctly: %s", err)
	}

	SetPort(db.MappedPort)
	e := echo.New()

	req := httptest.NewRequest(http.MethodPost, "/tenant/customer1/product/:pinstance/license", strings.NewReader(""))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded") //Q: are there no constants for this in go?
	rec := httptest.NewRecorder()

	echoCtx := e.NewContext(req, rec)
	echoCtx.SetPath("/tenant/:tenant/product/:pinstance/license")
	echoCtx.SetParamNames("tenant", "pinstance")
	echoCtx.SetParamValues("customer1", "p1")

	emptyBodyResp := GrantLicenseIfNotFull(echoCtx)

	if assert.NotNil(t, emptyBodyResp) {
		he, ok := emptyBodyResp.(*echo.HTTPError)
		if ok {
			assert.Equal(t, http.StatusBadRequest, he.Code)
			assert.Contains(t, he.Message, "Bad Request. User to grant access to needed")
		}
	}
}

/*
*
Get Licenses
*/
func TestGetLicenseReturnsListOfLicensedUsersForTenant(t *testing.T) {
	ctx := context.Background()
	/*using one tc instance per test bc don't quite know how to create fixtures and stuff. concious technical debt for now. enlighten me ;)*/
	db, err := setupSpiceDb(ctx, t)
	if err != nil {
		t.Fatalf("container not setup correctly: %s", err)
	}

	SetPort(db.MappedPort)
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/tenant/customer1/product/:pinstance/license?callingName=owner1", nil)
	rec := httptest.NewRecorder()

	echoCtx := e.NewContext(req, rec)
	echoCtx.SetPath("/tenant/:tenant/product/:pinstance/license")
	echoCtx.SetParamNames("tenant", "pinstance")
	echoCtx.SetParamValues("customer1", "p1")

	if assert.NoError(t, GetLicenseInfoForProductInstance(echoCtx)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		ja := jsonassert.New(t)
		// find some sort of payload
		name := "p1"
		active := 4 //relations for owner1, user1, user2 and user3
		max := 5
		ja.Assertf(rec.Body.String(), `
	{
		"name": "%s",
		"active_licenses": %d,
		"max_seats": %d
	}`, name, active, max)
	}
}

func TestGetLicenseForbiddenForOtherTenant(t *testing.T) {
	ctx := context.Background()
	/*using one tc instance per test bc don't quite know how to create fixtures and stuff. concious technical debt for now. enlighten me ;)*/
	db, err := setupSpiceDb(ctx, t)
	if err != nil {
		t.Fatalf("container not setup correctly: %s", err)
	}

	SetPort(db.MappedPort)
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/tenant/customer2/product/:pinstance/license?callingName=owner1", nil)
	rec := httptest.NewRecorder()

	echoCtx := e.NewContext(req, rec)
	echoCtx.SetPath("/tenant/:tenant/product/:pinstance/license")
	echoCtx.SetParamNames("tenant", "pinstance")
	echoCtx.SetParamValues("customer2", "p2")

	//weird weird way of testing for errors in echo
	err2 := GetLicenseInfoForProductInstance(echoCtx)
	if assert.NotNil(t, err2) {
		he, ok := err2.(*echo.HTTPError)
		if ok {
			assert.Equal(t, http.StatusForbidden, he.Code)
			assert.Contains(t, he.Message, "You are not allowed to see licensing information")
		}
	}
}

func TestGetLicenseForbiddenWithoutRightPermission(t *testing.T) {
	ctx := context.Background()
	/*using one tc instance per test bc don't quite know how to create fixtures and stuff. concious technical debt for now. enlighten me ;)*/
	db, err := setupSpiceDb(ctx, t)
	if err != nil {
		t.Fatalf("container not setup correctly: %s", err)
	}

	SetPort(db.MappedPort)
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/tenant/customer2/product/:pinstance/license?callingName=user1", nil)
	rec := httptest.NewRecorder()

	echoCtx := e.NewContext(req, rec)
	echoCtx.SetPath("/tenant/:tenant/product/:pinstance/license")
	echoCtx.SetParamNames("tenant", "pinstance")
	echoCtx.SetParamValues("customer2", "p2")

	//weird weird way of testing for errors in echo
	err2 := GetLicenseInfoForProductInstance(echoCtx)
	if assert.NotNil(t, err2) {
		he, ok := err2.(*echo.HTTPError)
		if ok {
			assert.Equal(t, http.StatusForbidden, he.Code)
			assert.Contains(t, he.Message, "You are not allowed to see licensing information")
		}
	}
}
