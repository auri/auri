package models

import (
	"auri/helpers"
	"time"

	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
)

//Reset is used by pop to map your .model.Name.Proper.Pluralize.Underscore database table to your go code.
type Reset struct {
	ID         uuid.UUID `db:"id"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
	Email      string    `db:"email" form:"mail"`
	Login      string    `db:"login"`
	Token      string    `db:"token"`
	ExpiryDate time.Time `db:"expiry_date"`
}

//ValidateEmail this will validate the email for password reset
func (reset *Reset) ValidateEmail(_ *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&emailValidator{Name: "email", Field: reset.Email}), nil
}

// NewToken generates and sets a new token
func (reset *Reset) NewToken() (err error) {
	reset.Token, err = helpers.NewUUID()
	return
}
