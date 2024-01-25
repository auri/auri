package popmw

import (
	"errors"
	"time"

	"github.com/gobuffalo/buffalo"
	pp "github.com/gobuffalo/buffalo-pop/v3/pop"
	"github.com/gobuffalo/events"
	"github.com/gobuffalo/pop/v6"
)

// PopTransaction is a piece of Buffalo middleware that wraps each
// request in a transaction. The transaction will automatically get
// committed if there's no errors and the response status code is a
// 2xx or 3xx, otherwise it'll be rolled back. It will also add a
// field to the log, "db", that shows the total duration spent during
// the request making database calls.
func Transaction(db *pop.Connection) buffalo.MiddlewareFunc {
	events.NamedListen("popmw.Transaction", func(e events.Event) {
		if e.Kind != "buffalo:app:start" {
			return
		}
		i, err := e.Payload.Pluck("app")
		if err != nil {
			return
		}
		if app, ok := i.(*buffalo.App); ok {
			pop.SetTxLogger(pp.TxLogger(app))
		}
	})

	return func(handler buffalo.Handler) buffalo.Handler {
		return func(c buffalo.Context) error {
			// wrap all requests in a transaction and set the length
			// of time doing things in the db to the log.
			// ANY error returned by the tx function will cause the
			// tx to be rolled back
			couldBeDBorYourErr := db.Transaction(func(tx *pop.Connection) error {
				// setup logging
				start := tx.Elapsed
				defer func() {
					finished := tx.Elapsed
					elapsed := time.Duration(finished - start)
					c.LogField("db", elapsed)
				}()

				c.LogField("conn", tx.ID)
				if tx.TX != nil {
					c.LogField("tx", tx.TX.ID)
				}

				// add the transaction to the context
				c.Set("tx", tx)

				// call the next handler; if it errors stop and return the error
				if yourError := handler(c); yourError != nil {
					return yourError
				}

				// check the response status code. if the code is NOT 200..399
				// then it is considered "NOT SUCCESSFUL" and an error will be returned
				if res, ok := c.Response().(*buffalo.Response); ok {
					if res.Status >= 400 {
						return errNonSuccess
					}
				}

				// as far was we can tell everything went well
				return nil
			})

			// couldBeDBorYourErr could be one of possible values:
			// * nil - everything went well, if so, return
			// * yourError - an error returned from your application, middleware, etc...
			// * a database error - this is returned if there were problems committing the transaction
			// * a errNonSuccess - this is returned if the response status code is 4xx or 5xx
			if couldBeDBorYourErr != nil && !errors.Is(couldBeDBorYourErr, errNonSuccess) {
				return couldBeDBorYourErr
			}
			return nil
		}
	}
}

var errNonSuccess = errors.New("non success status code")
