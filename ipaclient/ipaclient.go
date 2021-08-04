package ipaclient

import (
	"auri/config"
	"crypto/tls"
	"net/http"

	"github.com/tehwalris/go-freeipa/freeipa"
)

//IPA interface describes methods we use from IPA API via freeipa lib
type IPA interface {
	UserAdd(reqArgs *freeipa.UserAddArgs, optArgs *freeipa.UserAddOptionalArgs) (*freeipa.UserAddResult, error)
	Passwd(reqArgs *freeipa.PasswdArgs, optArgs *freeipa.PasswdOptionalArgs) (*freeipa.PasswdResult, error)
	UserFind(criteria string, reqArgs *freeipa.UserFindArgs, optArgs *freeipa.UserFindOptionalArgs) (*freeipa.UserFindResult, error)
	UserMod(reqArgs *freeipa.UserModArgs, optArgs *freeipa.UserModOptionalArgs) (*freeipa.UserModResult, error)
	UserShow(reqArgs *freeipa.UserShowArgs, optArgs *freeipa.UserShowOptionalArgs) (*freeipa.UserShowResult, error)
}

//ipa contains the initialized IPA client to work with IPA API
var ipa IPA

// New links to the constructor, which is used to create the client
var New = NewClient

//NewClient returns initialized IPA Client
func NewClient(host, user, password string, tlsSkipVerify bool) (IPA, error) {
	ipa, err := freeipa.Connect(
		host,
		&http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: tlsSkipVerify,
			},
		},
		user,
		password,
	)
	if err != nil {
		return nil, err
	}

	return ipa, nil
}

// GetClient returns the singleton instance of initialized IPA client
// GetClient is usually the proper way to get the initialized client instead of NewClient
func GetClient() (IPA, error) {
	var err error

	if ipa == nil {
		conf := config.GetInstance()

		// we store ipa client to the local temp var
		// as the function call alone sets type (see reflect.Type) to the interface
		// and our nil check above doesn't work as expected
		// https://groups.google.com/forum/#!msg/golang-nuts/wnH302gBa4I/i_cPI8Gu9P8J
		lipa, err := New(
			conf.IPAHost,
			conf.IPAUser,
			conf.IPAPassword,
			conf.IPATLSSkipVerify,
		)
		if err != nil {
			return nil, err
		}

		ipa = lipa
	}

	return ipa, err
}
