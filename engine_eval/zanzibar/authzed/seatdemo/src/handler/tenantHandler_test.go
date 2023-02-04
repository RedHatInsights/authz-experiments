package handler

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetTenantUserAccessReturnsListOfTenantUsers(t *testing.T) {
	ctx := context.Background()
	/*using one tc instance per test bc don't quite know how to create fixtures and stuff. concious technical debt for now. enlighten me ;)*/
	db, err := setupSpiceDb(ctx, t)
	if err != nil {
		t.Fatalf("container not setup correctly: %s", err)
	}

	SetPort(db.MappedPort)
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/tenant/customer1/user?callingName=owner1", nil)
	rec := httptest.NewRecorder()

	echoCtx := e.NewContext(req, rec)
	echoCtx.SetPath("/tenant/:tenant/user")
	echoCtx.SetParamNames("tenant")
	echoCtx.SetParamValues("customer1")

	if assert.NoError(t, GetTenantUsers(echoCtx)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "owner1")
		assert.Contains(t, rec.Body.String(), "user1")
		assert.Contains(t, rec.Body.String(), "user4")

	}
}

func TestOwnerOfOneTenantCannotAccessOtherTenantUserList(t *testing.T) {
	ctx := context.Background()
	/*using one tc instance per test bc don't quite know how to create fixtures and stuff. concious technical debt for now. enlighten me ;)*/
	db, err := setupSpiceDb(ctx, t)
	if err != nil {
		t.Fatalf("container not setup correctly: %s", err)
	}

	SetPort(db.MappedPort)
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/tenant/customer1/user?callingName=t2owner1", nil)
	rec := httptest.NewRecorder()

	echoCtx := e.NewContext(req, rec)
	echoCtx.SetPath("/tenant/:tenant/user")
	echoCtx.SetParamNames("tenant")
	echoCtx.SetParamValues("customer1")

	//weird weird way of testing for errors in echo
	err2 := GetTenantUsers(echoCtx)
	if assert.NotNil(t, err2) {
		he, ok := err2.(*echo.HTTPError)
		if ok {
			assert.Equal(t, http.StatusForbidden, he.Code)
			assert.Contains(t, he.Message, "nothing to see here")
		}
	}
}

func TestGetUsersNotAvailableForNormalUsers(t *testing.T) {
	ctx := context.Background()
	/*using one tc instance per test bc don't quite know how to create fixtures and stuff. concious technical debt for now. enlighten me ;)*/
	db, err := setupSpiceDb(ctx, t)
	if err != nil {
		t.Fatalf("container not setup correctly: %s", err)
	}

	SetPort(db.MappedPort)
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/tenant/customer1/user?callingName=user1", nil)
	rec := httptest.NewRecorder()

	echoCtx := e.NewContext(req, rec)
	echoCtx.SetPath("/tenant/:tenant/user")
	echoCtx.SetParamNames("tenant")
	echoCtx.SetParamValues("customer1")

	//weird weird way of testing for errors in echo
	err2 := GetTenantUsers(echoCtx)
	if assert.NotNil(t, err2) {
		he, ok := err2.(*echo.HTTPError)
		if ok {
			assert.Equal(t, http.StatusForbidden, he.Code)
			assert.Contains(t, he.Message, "nothing to see here")
		}
	}
}
