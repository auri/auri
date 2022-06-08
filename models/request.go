package models

import (
	"auri/config"
	"auri/helpers"
	"auri/ipaclient"
	"strconv"

	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

// Request represents the account request
type Request struct {
	ID                   uuid.UUID    `json:"id" db:"id"`
	CreatedAt            time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time    `json:"updated_at" db:"updated_at"`
	Email                string       `db:"email" form:"mail"`
	Name                 string       `db:"name" form:"name"`
	LastName             string       `db:"last_name" form:"lastname"`
	PublicKey            string       `db:"-" form:"ssh"`
	MailVerification     bool         `db:"email_verification"`
	ExpiryDate           time.Time    `db:"expiry_date"`
	Token                string       `json:"token" db:"token"`
	CommentField         nulls.String `db:"comment_field" form:"comment"`
	Password             string       `db:"-" form:"password"`
	PasswordConfirmation string       `db:"-" form:"passwordConfirmation"`
}

// Requests is a slice of Request type
type Requests []Request

//emailValidator used to validate email
type emailValidator struct {
	Field   string
	Name    string
	Message string
}

type passwordsshValidator struct {
	User    Request
	Message string
}

func (v *passwordsshValidator) IsValid(verrs *validate.Errors) {

	if len(v.User.PublicKey) != 0 {
		if err := ValidateSSH(v.User.PublicKey); err != nil {
			verrs.Add(validators.GenerateKey(v.User.PublicKey), err.Error())
			return
		}
	}
	if len(v.User.Password) != 0 {
		if err := VerifyPassword(v.User.Password); err != nil {
			verrs.Add(validators.GenerateKey(v.User.Password), err.Error())
			return
		}
	}

	if len(v.User.Password) == 0 && len(v.User.PublicKey) == 0 {
		verrs.Add(validators.GenerateKey(v.Message), "Please enter your ssh key or your password")
		return
	}

}

// IsValid will execute the function, which was passed as parameter,
// and will add an error to the list of errors if an error was occured
func (f *emailValidator) IsValid(verrs *validate.Errors) {
	//Split E-Mail string by @
	parts := strings.Split(f.Field, "@")
	if len(parts) != 2 {
		verrs.Add(validators.GenerateKey(f.Name), "Invalid email format")
	}
	if len(parts) == 2 {
		//to load up all contents in from json file
		hasErr := true

		for _, v := range config.GetInstance().AllowedDomains {
			if parts[1] == v {
				hasErr = false
				continue
			}
		}
		if hasErr {
			verrs.Add(validators.GenerateKey(f.Name), "Email does not match one of our permitted domains")
		}
	}
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (u *Request) Validate(tx *pop.Connection) (*validate.Errors, error) {
	var err error
	return validate.Validate(
		&validators.StringIsPresent{Field: u.Name, Name: "name", Message: "Name is required"},
		&validators.StringIsPresent{Field: u.LastName, Name: "lastname", Message: "Last name is required"},
		&emailValidator{Field: u.Email, Name: "email"},
		&validators.FuncValidator{
			Field:   u.Email,
			Name:    "email",
			Message: "Account with email %s already exists",
			Fn: func() bool {
				var b bool
				if u.Email != "" {
					b, err = ipaclient.EMailExists(u.Email)
					return !b
				}
				return true // TODO: this should be false, but to avoid multiple validation messages we return true here
			},
		},
		&validators.FuncValidator{
			Field:   u.Email,
			Name:    "email",
			Message: "Account request with email %s already exists",
			Fn: func() bool {
				var b bool
				q := tx.Where("email = ?", u.Email)
				if u.ID != uuid.Nil {
					q = q.Where("id != ?", u.ID)
				}
				if b, err = q.Exists(u); err != nil {
					err = errors.WithMessage(err, "DB error")
					return false
				}
				return !b
			},
		},
	), err
}

//ValidateSSHOrPassword this will validate the ssh key or the password
func (u *Request) ValidateSSHOrPassword() *validate.Errors {
	return validate.Validate(
		&validators.StringsMatch{Name: "password", Field: u.Password, Field2: u.PasswordConfirmation, Message: "Password and it's confirmation are not equal"},
		&passwordsshValidator{User: *u})
}

//Create gets run every time you call "pop.ValidateAndUpdate" method.
func (u *Request) Create(tx *pop.Connection) (*validate.Errors, error) {
	u.Email = strings.ToLower(strings.TrimSpace(u.Email))
	u.Name = strings.TrimSpace(u.Name)
	u.LastName = strings.TrimSpace(u.LastName)
	u.MailVerification = false
	u.ExpiryDate = helpers.TimeNowAddHours(config.GetInstance().EmailValidationTimeframe)

	if err := u.NewToken(); err != nil {
		return nil, err
	}

	return tx.ValidateAndCreate(u)
}

//PasswordHash returns the hashed password of user
func (u *Request) PasswordHash() (string, error) {
	//TODO: maybe we should raise an error here, let's return invalid hash for now
	if len(u.Password) == 0 {
		return "*", nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.MinCost)
	if err != nil {
		err = errors.WithMessage(err, "Can not create password hash")
	}
	return string(hash), nil
}

// NewToken generates and sets a new token
func (u *Request) NewToken() (err error) {
	u.Token, err = helpers.NewUUID()
	return
}

// Login generates a new login
// if index isn't zero it's added to the login
func (u *Request) Login(index int) string {
	var login string

	name := []rune(u.Name)
	if len(name) > 0 {
		login = login + string(name[0])
	}

	if len(u.LastName) > 0 {
		login = login + u.LastName
	}

	// replace german umlauts
	umlauts := map[string]string{
		"Ä": "Ae",
		"ä": "ae",
		"Ö": "Oe",
		"ö": "oe",
		"Ü": "Ue",
		"ü": "ue",
		"ß": "ss",
	}
	for k, v := range umlauts {
		login = strings.ReplaceAll(login, k, v)
	}

	//login isn't allowed to have spaces
	login = strings.ReplaceAll(login, " ", "")

	if index > 0 {
		login = login + strconv.Itoa(index)
	}
	return strings.ToLower(login)
}
