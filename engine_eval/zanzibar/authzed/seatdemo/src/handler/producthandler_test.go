package handler

import (
	"context"
	v1 "github.com/authzed/authzed-go/proto/authzed/api/v1"
	"github.com/kinbiko/jsonassert"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type EchoTestParams struct {
	methodType      string
	uri             string
	path            string
	bodyContent     string
	paramNames      []string
	paramValues     []string
	optionalHeaders map[string]string
}

func TestGrantLicenseReturnsBadRequestWhenNoProductInstanceLicenseFound(t *testing.T) {
	ctx := context.Background()
	/*using one tc instance per test bc don't quite know how to create fixtures and stuff. concious technical debt for now. enlighten me ;)*/
	db := getSpiceDbContainer(t, ctx)

	SetPort(db.MappedPort)

	p := EchoTestParams{
		methodType:      http.MethodPost,
		uri:             "/tenant/customer1/product/p9999/license",
		path:            "/tenant/:tenant/product/:pinstance/license",
		bodyContent:     "userId=user5",
		paramNames:      []string{"tenant", "pinstance"},
		paramValues:     []string{"customer1", "p9999"},
		optionalHeaders: map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
	}

	echoCtx, _ := populateEchoContext(p)

	errResp := GrantLicenseIfNotFull(echoCtx)

	assertHttpErrCodeAndMsg(t, http.StatusBadRequest, "No license found for product instance p9999", errResp)
}

func TestGrantLicenseGrantsLicenseIfAllConditionsMet(t *testing.T) {
	ctx := context.Background()
	/*using one tc instance per test bc don't quite know how to create fixtures and stuff. concious technical debt for now. enlighten me ;)*/
	db := getSpiceDbContainer(t, ctx)

	SetPort(db.MappedPort)

	p := EchoTestParams{
		methodType:      http.MethodPost,
		uri:             "/tenant/customer1/product/p1/license",
		path:            "/tenant/:tenant/product/:pinstance/license",
		bodyContent:     "userId=user5",
		paramNames:      []string{"tenant", "pinstance"},
		paramValues:     []string{"customer1", "p1"},
		optionalHeaders: map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
	}

	echoCtx, rec := populateEchoContext(p)

	if assert.NoError(t, GrantLicenseIfNotFull(echoCtx)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		ja := jsonassert.New(t)
		// find some sort of payload
		granted := true
		reason := "Successfully granted a license for instance p1 to user5. Remaining: 0"
		ja.Assertf(rec.Body.String(), `
	{
		"granted": %t,
		"reason": "%s"
	}`, granted, reason)
	}

	//cleanup bc container reuse.. TODO refactor
	client, _ := getSpiceDbApiClient(db.MappedPort)
	client.DeleteRelationships(ctx, &v1.DeleteRelationshipsRequest{
		RelationshipFilter: &v1.RelationshipFilter{
			ResourceType:       "user",
			OptionalResourceId: "user5",
			OptionalRelation:   "licensed_wsdm_user",
		},
	})

	client.DeleteRelationships(ctx, &v1.DeleteRelationshipsRequest{
		RelationshipFilter: &v1.RelationshipFilter{
			ResourceType:       "product_instance",
			OptionalResourceId: "p1",
			OptionalRelation:   "wsdm_user",
			OptionalSubjectFilter: &v1.SubjectFilter{
				SubjectType:       "user",
				OptionalSubjectId: "user5",
			},
		},
	})
}

func TestGrantLicenseReturns409ForUserWithAlreadyActivatedLicense(t *testing.T) {
	ctx := context.Background()
	/*using one tc instance per test bc don't quite know how to create fixtures and stuff. concious technical debt for now. enlighten me ;)*/
	db := getSpiceDbContainer(t, ctx)

	SetPort(db.MappedPort)
	p := EchoTestParams{
		methodType:      http.MethodPost,
		uri:             "/tenant/customer1/product/p1/license",
		path:            "/tenant/:tenant/product/:pinstance/license",
		bodyContent:     "userId=user1",
		paramNames:      []string{"tenant", "pinstance"},
		paramValues:     []string{"customer1", "p1"},
		optionalHeaders: map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
	}

	echoCtx, _ := populateEchoContext(p)
	errResp := GrantLicenseIfNotFull(echoCtx)

	assertHttpErrCodeAndMsg(t, http.StatusConflict, "Already active license for user user1 found.", errResp)
}

func TestGrantLicenseReturns403ForUserNotMemberOfTenant(t *testing.T) {
	ctx := context.Background()
	/*using one tc instance per test bc don't quite know how to create fixtures and stuff. concious technical debt for now. enlighten me ;)*/
	db := getSpiceDbContainer(t, ctx)

	SetPort(db.MappedPort)
	p := EchoTestParams{
		methodType:      http.MethodPost,
		uri:             "/tenant/customer1/product/p1/license",
		path:            "/tenant/:tenant/product/:pinstance/license",
		bodyContent:     "userId=t2user3",
		paramNames:      []string{"tenant", "pinstance"},
		paramValues:     []string{"customer1", "p1"},
		optionalHeaders: map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
	}

	echoCtx, _ := populateEchoContext(p)

	errResp := GrantLicenseIfNotFull(echoCtx)

	assertHttpErrCodeAndMsg(t, http.StatusForbidden, "User t2user3 is not a member of licensed tenant customer1", errResp)
}

func TestGrantLicenseRevokesGrantIfMaxReached(t *testing.T) {
	ctx := context.Background()
	/*using one tc instance per test bc don't quite know how to create fixtures and stuff. concious technical debt for now. enlighten me ;)*/
	db := getSpiceDbContainer(t, ctx)

	SetPort(db.MappedPort)
	p := EchoTestParams{
		methodType:      http.MethodPost,
		uri:             "/tenant/customer2/product/p2/license",
		path:            "/tenant/:tenant/product/:pinstance/license",
		bodyContent:     "userId=t2user3",
		paramNames:      []string{"tenant", "pinstance"},
		paramValues:     []string{"customer2", "p2"},
		optionalHeaders: map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
	}

	echoCtx, rec := populateEchoContext(p)

	if assert.NoError(t, GrantLicenseIfNotFull(echoCtx)) {
		assert.Equal(t, http.StatusConflict, rec.Code)
		ja := jsonassert.New(t)
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
	db := getSpiceDbContainer(t, ctx)

	SetPort(db.MappedPort)
	p := EchoTestParams{
		methodType:      http.MethodPost,
		uri:             "/tenant/customer1/product/p1/license",
		path:            "/tenant/:tenant/product/:pinstance/license",
		paramNames:      []string{"tenant", "pinstance"},
		paramValues:     []string{"customer2", "p2"},
		optionalHeaders: map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
	}

	echoCtx, _ := populateEchoContext(p)
	errResp := GrantLicenseIfNotFull(echoCtx)

	assertHttpErrCodeAndMsg(t, http.StatusBadRequest, "Bad Request. User to grant access to needed", errResp)
}

/*
*
Get Licenses
*/
func TestGetLicenseReturnsBadRequestIfNoLicenseFound(t *testing.T) {
	ctx := context.Background()
	/*using one tc instance per test bc don't quite know how to create fixtures and stuff. concious technical debt for now. enlighten me ;)*/
	db := getSpiceDbContainer(t, ctx)

	SetPort(db.MappedPort)
	p := EchoTestParams{
		methodType:  http.MethodGet,
		uri:         "/tenant/customer2/product/p999/license?callingName=t2owner",
		path:        "/tenant/:tenant/product/:pinstance/license",
		paramNames:  []string{"tenant", "pinstance"},
		paramValues: []string{"customer2", "p999"},
	}

	echoCtx, _ := populateEchoContext(p)
	errResp := GetLicenseInfoForProductInstance(echoCtx)

	assertHttpErrCodeAndMsg(t, http.StatusBadRequest, "No license found for product instance p999", errResp)
}

func TestGetLicenseReturnsListOfLicensedUsersForTenant(t *testing.T) {
	ctx := context.Background()
	/*using one tc instance per test bc don't quite know how to create fixtures and stuff. concious technical debt for now. enlighten me ;)*/
	db := getSpiceDbContainer(t, ctx)

	SetPort(db.MappedPort)

	p := EchoTestParams{
		methodType:  http.MethodGet,
		uri:         "/tenant/customer1/product/p1/license?callingName=owner1",
		path:        "/tenant/:tenant/product/:pinstance/license",
		paramNames:  []string{"tenant", "pinstance"},
		paramValues: []string{"customer1", "p1"},
	}

	echoCtx, rec := populateEchoContext(p)

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
	db := getSpiceDbContainer(t, ctx)

	SetPort(db.MappedPort)
	p := EchoTestParams{
		methodType:  http.MethodGet,
		uri:         "/tenant/customer2/product/p2/license?callingName=owner1",
		path:        "/tenant/:tenant/product/:pinstance/license",
		paramNames:  []string{"tenant", "pinstance"},
		paramValues: []string{"customer2", "p2"},
	}

	echoCtx, _ := populateEchoContext(p)

	errResp := GetLicenseInfoForProductInstance(echoCtx)

	assertHttpErrCodeAndMsg(t, http.StatusForbidden, "You are not allowed to see licensing information", errResp)
}

func TestGetLicenseForbiddenWithoutRightPermission(t *testing.T) {
	ctx := context.Background()
	/*using one tc instance per test bc don't quite know how to create fixtures and stuff. concious technical debt for now. enlighten me ;)*/
	db := getSpiceDbContainer(t, ctx)

	SetPort(db.MappedPort)
	p := EchoTestParams{
		methodType:  http.MethodGet,
		uri:         "/tenant/customer2/product/p2/license?callingName=t2user2",
		path:        "/tenant/:tenant/product/:pinstance/license",
		paramNames:  []string{"tenant", "pinstance"},
		paramValues: []string{"customer2", "p2"},
	}

	echoCtx, _ := populateEchoContext(p)

	errResp := GetLicenseInfoForProductInstance(echoCtx)
	assertHttpErrCodeAndMsg(t, http.StatusForbidden, "You are not allowed to see licensing information", errResp)
}

func populateEchoContext(p EchoTestParams) (echo.Context, *httptest.ResponseRecorder) {

	e := echo.New()
	req := httptest.NewRequest(p.methodType, p.uri, strings.NewReader(p.bodyContent))

	for k, v := range p.optionalHeaders {
		req.Header.Set(k, v)
	}

	rec := httptest.NewRecorder()

	echoCtx := e.NewContext(req, rec)
	echoCtx.SetPath(p.path)
	echoCtx.SetParamNames(p.paramNames...)
	echoCtx.SetParamValues(p.paramValues...) //p3 does not exist
	return echoCtx, rec
}

func assertHttpErrCodeAndMsg(t *testing.T, statusCode int, message string, err2 error) {
	if assert.NotNil(t, err2) {
		he, ok := err2.(*echo.HTTPError)
		if ok {
			assert.Equal(t, statusCode, he.Code)
			assert.Contains(t, he.Message, message)
		}
	}
}

func getSpiceDbContainer(t *testing.T, ctx context.Context) *spicedbContainer {
	db, err := setupSpiceDb(ctx, t)

	if err != nil {
		t.Fatalf("container not setup correctly: %s", err)
	}
	return db
}
