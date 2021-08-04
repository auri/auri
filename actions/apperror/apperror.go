package apperror

import (
	"github.com/gobuffalo/buffalo"
)

//ErrorType defines the type of error
type ErrorType uint8

const (
	//IPAError defines errors related to IPA
	IPAError ErrorType = iota
	//DBError defines errors related to database
	DBError
	//MailError defines errors related to the mail dispatching
	MailError
	//AuriError defines internal application error
	AuriError
)

//errorPrefixes sets up the specifics of different mails with tokens
var errorPrefixes = map[ErrorType]string{
	IPAError:  "IPA",
	DBError:   "Database",
	MailError: "Mail dispatching",
	AuriError: "Auri internal",
}

//InvokeError logs the error and displays the generic error message to the user
func InvokeError(c buffalo.Context, errType ErrorType, detErr error) {
	c.Logger().Errorf("%v error: %s", errorPrefixes[errType], detErr)
	c.Flash().Add("danger", "Something went wrong, please contact the system administrator")
}
