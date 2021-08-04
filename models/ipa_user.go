package models

import (
	"auri/helpers"
	"time"

	"github.com/gobuffalo/nulls"
	"github.com/gofrs/uuid"
)

//IpaUser this table used to get user from IPA
type IpaUser struct {
	ID                uuid.UUID  `db:"id"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at" db:"updated_at"`
	UID               string     `db:"uid"`
	RemindedAt        time.Time  `db:"reminded_at"`
	Token             string     `db:"token"`
	NotifiedAt        nulls.Time `db:"notified_at"`
	NotificationsSent int        `db:"notifications_sent"`
}

// NewToken generates and sets a new token
func (u *IpaUser) NewToken() (err error) {
	u.Token, err = helpers.NewUUID()
	return
}
