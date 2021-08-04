package ipaclient

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/tehwalris/go-freeipa/freeipa"
)

//AddUser add new user to IPA server
func AddUser(login, name, lastname, email, pubKey string) error {
	return runIPATransaction("Adding user with login '"+login+"' and email '"+email+"'", func(client IPA) error {
		_, err := client.UserAdd(&freeipa.UserAddArgs{
			Givenname: name,
			Sn:        lastname,
		}, &freeipa.UserAddOptionalArgs{
			UID:          &login,
			Mail:         &[]string{email},
			Ipasshpubkey: &[]string{pubKey},
		})
		if err != nil {
			return errors.WithMessage(err, "Can not add user to the server")
		}

		return err
	})
}

// LoginExists checks if given login exists on the server
func LoginExists(login string) (bool, error) {
	var err error
	loginExists := false
	err = runIPATransaction("Checking if login '"+login+"' exists", func(client IPA) error {
		var find *freeipa.UserFindResult
		find, err = client.UserFind("", &freeipa.UserFindArgs{}, &freeipa.UserFindOptionalArgs{
			UID: &login,
		})
		if err != nil {
			return errors.WithMessage(err, "Can not check if login exists")
		}

		if len(find.Result) > 1 {
			return errors.New("Got multiple results for lookup of login " + login)
		}

		if len(find.Result) > 0 {
			loginExists = true
		}

		return nil
	})

	return loginExists, err
}

// EMailExists checks if given mail exists on the server
func EMailExists(mail string) (bool, error) {
	var err error
	mailExists := false
	err = runIPATransaction("Checking if mail '"+mail+"' exists", func(client IPA) error {
		var find *freeipa.UserFindResult
		find, err = client.UserFind("", &freeipa.UserFindArgs{}, &freeipa.UserFindOptionalArgs{
			Mail: &[]string{mail},
		})
		if err != nil {
			return errors.WithMessage(err, "Can not check if mail exists")
		}

		if len(find.Result) > 0 {
			mailExists = true
		}

		return nil
	})

	return mailExists, err
}

//UserFind tries to find a login for a given email address in IPA
// returns nil if no login is found
func UserFind(email string) (string, error) {
	var login string
	var err error
	err = runIPATransaction("Searching IPA for a user with email address "+email, func(client IPA) error {
		var find *freeipa.UserFindResult
		find, err = client.UserFind("", &freeipa.UserFindArgs{}, &freeipa.UserFindOptionalArgs{
			Mail: &[]string{email},
		})
		if err != nil {
			return errors.WithMessage(err, "Can not find user on IPA server")
		}

		if len(find.Result) == 0 {
			return nil
		}

		if len(find.Result) != 1 {
			return errors.New("Multiple accounts found for email " + email)
		}

		// we don't want to rely completely on the IPA UserFind, as additional safety step we want to compare
		// emails via string matching to be 100% sure
		for _, ml := range *find.Result[0].Mail {
			if strings.ToLower(strings.TrimSpace(ml)) == strings.ToLower(strings.TrimSpace(email)) {
				login = find.Result[0].UID
			}
		}

		return err
	})
	return login, err
}

// GetGroups returns groups (direct and indirect) where a given user is a member
func GetGroups(login string) ([]string, error) {
	var memberOf []string
	var err error
	err = runIPATransaction("Getting membership information of user with login "+login, func(client IPA) error {
		var find *freeipa.UserShowResult
		find, err = client.UserShow(&freeipa.UserShowArgs{}, &freeipa.UserShowOptionalArgs{
			UID: &login,
		})
		if err != nil {
			return errors.WithMessage(err, "Can't get data of user")
		}

		if find.Result.MemberofGroup != nil {
			for _, gr := range *find.Result.MemberofGroup {
				memberOf = append(memberOf, gr)
			}
		}
		if find.Result.MemberofindirectGroup != nil {
			for _, gr := range *find.Result.MemberofindirectGroup {
				memberOf = append(memberOf, gr)
			}
		}
		return nil
	})
	return memberOf, err
}

//UserModify this function will modify user information on ipa
func UserModify(login, password, pubKey string) error {
	return runIPATransaction("Setting authentication credentials for user "+login, func(client IPA) error {
		var paramPass *string
		if len(password) > 0 {
			paramPass = &password
		}

		var paramKey *[]string
		if len(pubKey) > 0 {
			paramKey = &[]string{pubKey}
		}

		_, err := client.UserMod(&freeipa.UserModArgs{}, &freeipa.UserModOptionalArgs{
			UID:          &login,
			Userpassword: paramPass,
			Ipasshpubkey: paramKey,
		})
		if err != nil {
			return errors.WithMessage(err, "Can't modify user information on the server for user with login '"+login+"'")
		}

		// when we change a password, IPA automatically sets krbPasswordExpiration and forces the user to
		// change the password on the next IPA login. As we changed the password via Auri, it doesn't make really sense
		if paramPass != nil {
			_, err = client.UserMod(&freeipa.UserModArgs{}, &freeipa.UserModOptionalArgs{
				UID:     &login,
				Setattr: &[]string{"krbPasswordExpiration=del"},
				Delattr: &[]string{"krbPasswordExpiration=del"},
			})
			if err != nil {
				return errors.WithMessage(err, "Can't reset password expiration time on the server")
			}
		}

		return err
	})
}
