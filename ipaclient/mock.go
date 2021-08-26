package ipaclient

import (
	"strings"

	"github.com/pkg/errors"

	"github.com/tehwalris/go-freeipa/freeipa"
)

func sPointer(s string) *string {
	return &s
}

// ipaMock is used as a mock client for IPA
type ipaMock struct {
	userAdded    *freeipa.User
	userModified *freeipa.User
}

// UserAdd is mocking freeipa.Client.UserAdd()
func (moc *ipaMock) UserAdd(args *freeipa.UserAddArgs, oargs *freeipa.UserAddOptionalArgs) (*freeipa.UserAddResult, error) {
	uid := ""
	if oargs.UID != nil {
		uid = *oargs.UID
	}
	moc.userAdded = &freeipa.User{UID: uid,
		Givenname: &args.Givenname,
		Sn:        args.Sn,
		Mail:      oargs.Mail}

	return &freeipa.UserAddResult{Result: *moc.userAdded}, nil
}

// UserFind is mocking freeipa.Client.UserFind()
func (moc *ipaMock) UserFind(_ string, _ *freeipa.UserFindArgs, optionalArgs *freeipa.UserFindOptionalArgs) (*freeipa.UserFindResult, error) {
	var login, mail string
	if optionalArgs.UID != nil {
		login = *optionalArgs.UID
	}
	if optionalArgs.Mail != nil && len(*optionalArgs.Mail) > 0 {
		mail = strings.ToLower(strings.TrimSpace((*optionalArgs.Mail)[0]))
	}

	if login == "user-present-in-IPA" || mail == "user-present-in-ipa@example.com" {
		return &freeipa.UserFindResult{
			Result: []freeipa.User{
				{
					UID:              "user-present-in-IPA",
					Givenname:        sPointer("Test"),
					Sn:               "User",
					Mail:             &[]string{"user-present-in-IPA@example.com"},
					Krbcanonicalname: sPointer("user-present-in-IPA@example.com"),
				},
			},
		}, nil
	} else if login == "admin-user-present-in-IPA" || mail == "admin-user-present-in-ipa@example.com" {
		return &freeipa.UserFindResult{
			Result: []freeipa.User{
				{
					UID:              "admin-user-present-in-IPA",
					Givenname:        sPointer("Admin"),
					Sn:               "Admin User",
					Mail:             &[]string{"admin-user-present-in-IPA@example.com"},
					Krbcanonicalname: sPointer("admin-user-present-in-IPA@example.com"),
				},
			},
		}, nil
	} else if login == "" && mail == "" {
		return &freeipa.UserFindResult{
			Result: []freeipa.User{
				{
					UID:              "user-present-in-IPA",
					Givenname:        sPointer("Test"),
					Sn:               "User",
					Mail:             &[]string{"user-present-in-IPA@example.com"},
					Krbcanonicalname: sPointer("user-present-in-IPA@example.com"),
				},
				{
					UID:              "admin-user-present-in-IPA",
					Givenname:        sPointer("Admin"),
					Sn:               "Admin User",
					Mail:             &[]string{"admin-user-present-in-IPA@example.com"},
					Krbcanonicalname: sPointer("admin-user-present-in-IPA@example.com"),
				},
			},
		}, nil
	}

	return &freeipa.UserFindResult{}, nil
}

// UserMod is mocking freeipa.Client.UserMod()
func (moc *ipaMock) UserMod(_ *freeipa.UserModArgs, oargs *freeipa.UserModOptionalArgs) (*freeipa.UserModResult, error) {
	moc.userModified = &freeipa.User{UID: *oargs.UID,
		Userpassword: oargs.Userpassword,
		Ipasshpubkey: oargs.Ipasshpubkey}

	return &freeipa.UserModResult{Result: *moc.userModified}, nil
}

// Passwd is mocking freeipa.Client.Passwd()
func (moc *ipaMock) Passwd(_ *freeipa.PasswdArgs, _ *freeipa.PasswdOptionalArgs) (*freeipa.PasswdResult, error) {
	return &freeipa.PasswdResult{}, nil
}

func (moc *ipaMock) UserShow(_ *freeipa.UserShowArgs, optionalArgs *freeipa.UserShowOptionalArgs) (*freeipa.UserShowResult, error) {
	if optionalArgs.UID == nil {
		return &freeipa.UserShowResult{}, errors.New("Missing user login")
	}

	if *optionalArgs.UID == "user-present-in-IPA" {
		return &freeipa.UserShowResult{
			Result: freeipa.User{
				UID:                   "user-present-in-IPA",
				Givenname:             sPointer("Test"),
				Sn:                    "User",
				Mail:                  &[]string{"user-present-in-IPA@example.com"},
				Krbcanonicalname:      sPointer("user-present-in-IPA@example.com"),
				MemberofGroup:         &[]string{"group1", "group2"},
				MemberofindirectGroup: &[]string{"indirectgroup1", "indirectgroup2"},
			},
		}, nil
	} else if *optionalArgs.UID == "admin-user-present-in-IPA" {
		return &freeipa.UserShowResult{
			Result: freeipa.User{
				UID:                   "admin-user-present-in-IPA",
				Givenname:             sPointer("Admin"),
				Sn:                    "Admin User",
				Mail:                  &[]string{"admin-user-present-in-IPA@example.com"},
				Krbcanonicalname:      sPointer("admin-user-present-in-IPA@example.com"),
				MemberofGroup:         &[]string{"admins", "group2"},
				MemberofindirectGroup: &[]string{"indirectadmins", "indirectgroup2"},
			},
		}, nil
	} else {
		return &freeipa.UserShowResult{}, errors.Errorf("NotFound (4001): %v: user not found", *optionalArgs.UID)
	}
}

//NewMock returns initialized IPA Mock
func NewMock(_, _, _ string, _ bool) (IPA, error) {
	return &ipaMock{}, nil
}

// UserAddedFromMock returns the last created user from mock
func UserAddedFromMock() *freeipa.User {
	ipaclient, err := GetClient()
	if err != nil {
		panic(err)
	}
	moc, ok := ipaclient.(*ipaMock)
	if ok != true {
		panic("IPAClient is expected to be ipaMock")
	}

	return moc.userAdded
}

// UserModifiedFromMock returns the last modified user from mock
func UserModifiedFromMock() *freeipa.User {
	ipaclient, err := GetClient()
	if err != nil {
		panic(err)
	}
	moc, ok := ipaclient.(*ipaMock)
	if ok != true {
		panic("IPAClient is expected to be ipaMock")
	}

	return moc.userModified
}
