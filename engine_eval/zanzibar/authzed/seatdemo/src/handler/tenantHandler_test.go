package handler

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTenantUserAccessReturnsListOfTenantUsers(t *testing.T) {
	resp := runRequest(get("/tenant/customer1/user?callingName=owner1"))

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assertJsonResponse(t, resp, http.StatusOK, `["<<UNORDERED>>",{"userid":"user1"},{"userid":"user2"},{"userid":"user3"},{"userid":"user4"},{"userid":"user5"},{"userid":"owner1"}]`)
}

func TestOwnerOfOneTenantCannotAccessOtherTenantUserList(t *testing.T) {
	resp := runRequest(get("/tenant/customer1/user?callingName=t2owner1"))
	assertHttpErrCodeAndMsg(t, http.StatusForbidden, "nothing to see here", resp)
}

func TestGetUsersNotAvailableForNormalUsers(t *testing.T) {
	resp := runRequest(get("/tenant/customer1/user?callingName=user1"))
	assertHttpErrCodeAndMsg(t, http.StatusForbidden, "nothing to see here", resp)
}
