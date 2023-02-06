package handler

import (
	"context"
	"net/http"
	"testing"

	v1 "github.com/authzed/authzed-go/proto/authzed/api/v1"
	"github.com/stretchr/testify/assert"
)

func TestGrantLicenseReturnsBadRequestWhenNoProductInstanceLicenseFound(t *testing.T) {
	resp := runRequest(post("/tenant/customer1/product/p9999/license", "userId=user5"))
	assertHttpErrCodeAndMsg(t, http.StatusBadRequest, "No license found for product instance p9999", resp)
}

func TestGrantLicenseGrantsLicenseIfAllConditionsMet(t *testing.T) {
	ctx := context.Background()
	/*using one tc instance per test bc don't quite know how to create fixtures and stuff. concious technical debt for now. enlighten me ;)*/

	resp := runRequest(post("/tenant/customer1/product/p1/license", "userId=user5"))

	granted := true
	reason := "Successfully granted a license for instance p1 to user5. Remaining: 0"
	assertJsonResponse(t, resp, http.StatusCreated,
		`{
			"granted": %t,
			"reason": "%s"
		}`, granted, reason)

	//cleanup bc container reuse.. TODO refactor
	db, err := getSpiceDbContainer()
	assert.NoError(t, err)
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
	resp := runRequest(post("/tenant/customer1/product/p1/license", "userId=user1"))
	assertHttpErrCodeAndMsg(t, http.StatusConflict, "Already active license for user user1 found.", resp)
}

func TestGrantLicenseReturns403ForUserNotMemberOfTenant(t *testing.T) {
	resp := runRequest(post("/tenant/customer1/product/p1/license", "userId=t2user3"))
	assertHttpErrCodeAndMsg(t, http.StatusForbidden, "User t2user3 is not a member of licensed tenant customer1", resp)
}

func TestGrantLicenseRevokesGrantIfMaxReached(t *testing.T) {
	resp := runRequest(post("/tenant/customer2/product/p2/license", "userId=t2user3"))

	granted := false
	reason := "Maximum seats exceeded. Please extend your license."
	assertJsonResponse(t, resp, http.StatusConflict,
		`{
			"granted": %t,
			"reason": "%s"
		}`, granted, reason)
}

func TestGrantLicenseReturnsBadRequestWithoutBody(t *testing.T) {
	resp := runRequest(post("/tenant/customer1/product/p1/license", ""))
	assertHttpErrCodeAndMsg(t, http.StatusBadRequest, "Bad Request. User to grant access to needed", resp)
}

/*
*
Get Licenses
*/
func TestGetLicenseReturnsBadRequestIfNoLicenseFound(t *testing.T) {
	resp := runRequest(get("/tenant/customer2/product/p999/license?callingName=t2owner"))
	assertHttpErrCodeAndMsg(t, http.StatusBadRequest, "No license found for product instance p999", resp)
}

func TestGetLicenseReturnsListOfLicensedUsersForTenant(t *testing.T) {
	resp := runRequest(get("/tenant/customer1/product/p1/license?callingName=owner1"))

	name := "p1"
	active := 4 //relations for owner1, user1, user2 and user3
	max := 5
	assertJsonResponse(t, resp, http.StatusOK,
		`{
			"name": "%s",
			"active_licenses": %d,
			"max_seats": %d
		}`, name, active, max)
}

func TestGetLicenseForbiddenForOtherTenant(t *testing.T) {
	resp := runRequest(get("/tenant/customer2/product/p2/license?callingName=owner1"))
	assertHttpErrCodeAndMsg(t, http.StatusForbidden, "You are not allowed to see licensing information", resp)
}

func TestGetLicenseForbiddenWithoutRightPermission(t *testing.T) {
	resp := runRequest(get("/tenant/customer2/product/p2/license?callingName=t2user2"))
	assertHttpErrCodeAndMsg(t, http.StatusForbidden, "You are not allowed to see licensing information", resp)
}
