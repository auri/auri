package grifts

import (
	"auri/models"

	"github.com/markbates/grift/grift"
)

func (ts *TaskSuite) Test_CleanupRequestsExpiredTokenNotVerified() {
	ts.LoadFixture("request-unverified-expiredToken-not-in-ipa")
	request := models.Request{}
	ts.NoError(ts.DB.Eager().Find(&request, "f19d5ab3-e913-4c08-b4dd-627e26ccbc0b"))
	ctx := grift.Context{}
	ts.False(request.MailVerification)
	ts.NoError(grift.Run("cleanup_requests", &ctx))
	ts.Error(ts.DB.Eager().Find(&request, "f19d5ab3-e913-4c08-b4dd-627e26ccbc0b"))
}

func (ts *TaskSuite) Test_CleanupRequestsExpiredTokenVerified() {
	ts.LoadFixture("request-verified-expiredToken-not-in-ipa")
	request := models.Request{}
	ts.NoError(ts.DB.Eager().Find(&request, "f19d5ab3-e913-4c08-b4dd-627e26ccbc0b"))
	ctx := grift.Context{}
	ts.True(request.MailVerification)
	ts.NoError(grift.Run("cleanup_requests", &ctx))
	ts.NoError(ts.DB.Eager().Find(&request, "f19d5ab3-e913-4c08-b4dd-627e26ccbc0b"))
}

func (ts *TaskSuite) Test_CleanupRequestsValidTokenNotVerified() {
	ts.LoadFixture("request-unverified-validToken-not-in-ipa")
	request := models.Request{}
	ts.NoError(ts.DB.Eager().Find(&request, "f19d5ab3-e913-4c08-b4dd-627e26ccbc0b"))
	ctx := grift.Context{}
	ts.False(request.MailVerification)
	ts.NoError(grift.Run("cleanup_requests", &ctx))
	ts.NoError(ts.DB.Eager().Find(&request, "f19d5ab3-e913-4c08-b4dd-627e26ccbc0b"))
}
