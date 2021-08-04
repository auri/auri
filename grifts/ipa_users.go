// +build !prodBuild

package grifts

import (
	"auri/config"
	"auri/helpers"
	"auri/ipaclient"
	"auri/mail"
	"auri/models"
	"log"
	"strconv"

	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop/v5"
	"github.com/markbates/grift/grift"
	"github.com/pkg/errors"
	"github.com/tehwalris/go-freeipa/freeipa"
)

//This task will pull all user from IPA and will save them into
//shadow DB to send validation emails
var _ = grift.Desc("db:fetch_ipa_user", "This task will pull all user from ip and will save them into shadow DB to send validation email")
var _ = grift.Add("db:fetch_ipa_user", func(c *grift.Context) error {

	var err error
	ipa, err := ipaclient.GetClient()
	if err != nil {
		return errors.WithStack(err)
	}

	find, err := ipa.UserFind("", &freeipa.UserFindArgs{}, &freeipa.UserFindOptionalArgs{})
	if err != nil {
		return errors.WithStack(err)
	}

	if err := models.DB.Transaction(func(tx *pop.Connection) error {
		for _, getUser := range find.Result {
			user := models.IpaUser{UID: getUser.UID}
			uid := tx.Where("uid = ?", getUser.UID)
			exist, err := uid.Exists(user)
			if err != nil {
				return errors.Errorf("Database error %s ", err)
			}
			if exist {
				log.Printf("User: %s already exists in DB", getUser.UID)
				continue
			}

			if err := tx.Create(&user); err != nil {
				return errors.WithMessage(err, "Unable to add user to DB")

			}
		}
		return nil

	}); err != nil {
		return errors.WithMessage(err, "No transaction found")
	}
	return err
})

//This task will look in the DB for users to send an email to user
//to validate their account on IPA
var _ = grift.Desc("send_validation_email", "This task will look in DB for users to send an email to user to extend his account on IPA")
var _ = grift.Add("send_validation_email", func(c *grift.Context) error {

	var err error

	ipa, err := ipaclient.GetClient()
	if err != nil {
		return errors.WithStack(err)
	}

	conf := config.GetInstance()

	users := []models.IpaUser{}
	user := &models.IpaUser{}

	if err := models.DB.Transaction(func(tx *pop.Connection) error {

		if err := tx.RawQuery("select * from ipa_users where reminded_at < now()-interval '" + strconv.Itoa(conf.AccountCheckTimeframe) + "' month and notifications_sent = '0'").All(&users); err != nil {
			return errors.Errorf("Database error %s ", err)
		}

		if err := tx.RawQuery("select * from ipa_users where notified_at < now()-interval '" +
			strconv.Itoa(conf.AccountCheckReminderAfter) + "' day and notifications_sent <= '" +
			strconv.Itoa(conf.AccountCheckMaxReminders) + "'").All(&users); err != nil {
			return errors.Errorf("Database error %s ", err)
		}

		for _, sendEmailToUser := range users {

			if err := tx.Where("uid = ? ", sendEmailToUser.UID).First(user); err != nil {
				return errors.Errorf("Database error %s ", err)
			}

			getUserEmail, err := ipa.UserFind("", &freeipa.UserFindArgs{}, &freeipa.UserFindOptionalArgs{
				UID: &user.UID,
			})
			if err != nil {
				return errors.WithStack(err)
			}
			if user.Token == "" {
				if err := user.NewToken(); err != nil {
					return err
				}
			}
			user.NotifiedAt = nulls.NewTime(helpers.TimeNow())
			user.NotificationsSent++

			if err := tx.Update(user); err != nil {
				return errors.WithMessage(err, "User can't be updated")
			}

			emails := getUserEmail.Result[0].Mail
			if emails == nil || len(*emails) == 0 {
				log.Printf("No Email for user %s was found", getUserEmail.Result[0].UID)
				continue
			}
			if len(*emails) > 1 {
				log.Printf("Multiple emails found for user %s, using the first one", getUserEmail.Result[0].UID)
			}
			email := (*emails)[0]

			if err := mail.SendTokenMail(mail.TokenMailOpts{Email: email, Token: user.Token}, mail.AccountValidation); err != nil {
				return errors.WithMessage(err, "Something went wrong and can't send email to user")
			}
		}
		return err

	}); err != nil {
		return errors.WithMessage(err, "No transaction found")
	}
	return err
})
