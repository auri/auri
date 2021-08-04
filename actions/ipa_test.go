package actions

import (
	"auri/helpers"
	"auri/models"
	"time"

	"github.com/gobuffalo/nulls"
)

func (as *ActionSuite) Test_ExtendUserAccountTokenValidOK() {
	as.LoadFixture("extend_user_account")
	user := &models.IpaUser{}
	as.NoError(as.DB.Find(user, "f19d5ab3-e913-4c08-b4dd-627e26ccbc0b"))
	user.Token = "37e16abd-b69e-4b65-a3cd-a9c7b250f7a0"
	req := as.JSON("/user/ipa-account/validation/%s/", user.Token).Get()
	as.Equal(200, req.Code, req.Body.String())
	as.NoError(as.DB.Where("notifications_sent = 0 ").First(user))
	as.NoError(as.DB.Where("token = ? ", "").First(user))
	as.Equal(user.RemindedAt.Truncate(time.Second), helpers.TimeNow().Truncate(time.Second))
}
func (as *ActionSuite) Test_ExtendUserAccountTokenValidNO() {
	as.LoadFixture("extend_user_account")
	user := &models.IpaUser{}
	as.NoError(as.DB.Find(user, "f19d5ab3-e913-4c08-b4dd-627e26ccbc0b"))
	user.Token = "37e16abd-b69e-4b65-a3cd-a9c7b250f7a2"
	req := as.JSON("/user/ipa-account/validation/%s/", user.Token).Get()
	as.Equal(400, req.Code, req.Body.String())
}

func (as *ActionSuite) Test_ExtendUserAccountTokenExpired() {
	as.LoadFixture("extend_user_account")
	user := &models.IpaUser{}
	as.NoError(as.DB.Find(user, "f19d5ab3-e913-4c08-b4dd-627e26ccbc0b"))
	user.NotifiedAt = nulls.NewTime(helpers.TimeNowAddHours(49))
	as.NoError(as.DB.Eager().Save(user))
	user.Token = "37e16abd-b69e-4b65-a3cd-a9c7b250f7a0"
	req := as.JSON("/user/ipa-account/validation/%s/", user.Token).Get()
	as.Equal(400, req.Code, req.Body.String())
}
