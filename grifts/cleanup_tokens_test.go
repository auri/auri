package grifts

import (
	"auri/models"

	"github.com/markbates/grift/grift"
)

func (ts *TaskSuite) Test_CleanupTokensExpired() {
	ts.LoadFixture("resetToken-expired")
	resetToken := models.Reset{}
	ts.NoError(ts.DB.Eager().Find(&resetToken, "f19d5ab3-e913-4c08-b4dd-627e26ccbc0b"))
	ts.NoError(ts.DB.Eager().Where("expiry_date < now()").First(&resetToken))
	ctx := grift.Context{}
	ts.NoError(grift.Run("cleanup_reset_tokens", &ctx))
	ts.Error(ts.DB.Eager().Find(&resetToken, "f19d5ab3-e913-4c08-b4dd-627e26ccbc0b"))

}

func (ts *TaskSuite) Test_CleanupTokensNotExpired() {
	ts.LoadFixture("resetToken-valid")
	resetToken := models.Reset{}
	ts.NoError(ts.DB.Eager().Find(&resetToken, "f19d5ab3-e913-4c08-b4dd-627e26ccbc0b"))
	ctx := grift.Context{}
	ts.NoError(grift.Run("cleanup_reset_tokens", &ctx))
	ts.NoError(ts.DB.Eager().Find(&resetToken, "f19d5ab3-e913-4c08-b4dd-627e26ccbc0b"))
}
