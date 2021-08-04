// +build !prodBuild

package grifts

import (
	"auri/helpers"
	"auri/mail"
	"auri/models"

	"github.com/gobuffalo/nulls"
	"github.com/markbates/grift/grift"
)

func (ts *TaskSuite) Test_SaveUserInDB() {
	ctx := grift.Context{}
	ts.NoError(grift.Run("db:fetch_ipa_user", &ctx))
	user := models.IpaUser{UID: "user-present-in-IPA"}
	ts.NoError(ts.DB.Where("uid = ?", user.UID).First(&user))
}

func (ts *TaskSuite) Test_IPAUsersNoDuplicateEntries() {
	ts.LoadFixture("ipa_user_reminded_now_without_notifications")
	user := &models.IpaUser{}
	ts.NoError(ts.DB.Find(user, "f19d5ab3-e913-4c08-b4dd-627e26ccbc0b"))
	ctx := grift.Context{}
	ts.NoError(grift.Run("db:fetch_ipa_user", &ctx))
	user = &models.IpaUser{UID: "user-present-in-IPA"}
	query := ts.DB.Where("uid = ?", user.UID)
	userExist, err := query.Exists(user)
	ts.NoError(err)
	count, err := query.Count(user)
	ts.Equal(1, count)
	ts.NoError(err)
	ts.True(userExist)
}

func (ts *TaskSuite) Test_IPAUsersSendEmailForFirstTime() {
	ts.LoadFixture("ipa_user_reminded_1y_ago_without_notifications")
	user := &models.IpaUser{}
	ts.NoError(ts.DB.Find(user, "f19d5ab3-e913-4c08-b4dd-627e26ccbc0b"))
	ts.NoError(ts.DB.Where("notifications_sent = 0 ").First(user))
	ctx := grift.Context{}
	ts.NoError(grift.Run("send_validation_email", &ctx))
	ts.NoError(ts.DB.Find(user, "f19d5ab3-e913-4c08-b4dd-627e26ccbc0b"))
	ts.Equal(1, user.NotificationsSent)
	ts.NotEmpty(user.Token)
	ts.NoError(ts.DB.Where("notifications_sent = 1 ").First(user))
}

func (ts *TaskSuite) Test_Test_IPAUsersSendEmailForTheSecondTime() {
	ts.LoadFixture("send_email_no_reaction")
	user := &models.IpaUser{}
	ts.NoError(ts.DB.Find(user, "f19d5ab3-e913-4c08-b4dd-627e26ccbc0b"))
	ts.NoError(ts.DB.Where("notifications_sent = 1 ").First(user))
	ctx := grift.Context{}
	user.NotifiedAt = nulls.NewTime(helpers.TimeNowAddHours(-49))
	ts.NoError(ts.DB.Eager().Save(user))
	ts.NoError(grift.Run("send_validation_email", &ctx))
	ts.NoError(ts.DB.Find(user, "f19d5ab3-e913-4c08-b4dd-627e26ccbc0b"))
	ts.NoError(ts.DB.Where("notifications_sent = 2 ").First(user))
}

func (ts *TaskSuite) Test_SendEmailFirstTimeAndNoRreactionCaseSecondtime() {
	defer mail.ClearMock()
	ts.LoadFixture("send_email_mixed_users")

	user := &models.IpaUser{}
	ts.NoError(ts.DB.Find(user, "f19d5ab3-e913-4c08-b4dd-627e26ccbc0b"))
	ts.NoError(ts.DB.Where("notifications_sent = 0 ").First(user))
	ctx := grift.Context{}

	ts.NoError(grift.Run("send_validation_email", &ctx))
	ts.NoError(ts.DB.Find(user, "f19d5ab3-e913-4c08-b4dd-627e26ccbc0b"))
	ts.NoError(ts.DB.Where("notifications_sent = 1 ").First(user))

	ts.NoError(ts.DB.Find(user, "177c2e3c-8d19-40cb-bf1c-aabf45653aa0"))
	ts.NoError(ts.DB.Where("notifications_sent = 1 ").First(user))
	ts.NoError(grift.Run("send_validation_email", &ctx))
	ts.NoError(ts.DB.Find(user, "177c2e3c-8d19-40cb-bf1c-aabf45653aa0"))
	ts.NoError(ts.DB.Where("notifications_sent = 2 ").First(user))
}
