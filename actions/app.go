// Package actions provides the main gobuffalo application logic
package actions

import (
	"auri/config"
	"auri/logger"
	"auri/models"
	"auri/public"

	"log"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo-pop/v3/pop/popmw"
	csrf "github.com/gobuffalo/mw-csrf"
	paramlogger "github.com/gobuffalo/mw-paramlogger"
	"github.com/gorilla/sessions"
)

var app *buffalo.App

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
//
// Routing, middleware, groups, etc... are declared TOP -> DOWN.
// This means if you add a middleware to `app` *after* declaring a
// group, that group will NOT have that new middleware. The same
// is true of resource declarations as well.
//
// It also means that routes are checked in the order they are declared.
// `ServeFiles` is a CATCH-ALL route, so it should always be
// placed last in the route declarations, as it will prevent routes
// declared after it to never be called.
func App() *buffalo.App {
	conf := config.GetInstance()

	// DoInit for manual initialization
	DoInit()
	if app == nil {
		app = buffalo.New(buffalo.Options{
			Env:         conf.Environment,
			SessionName: "_auri_session",
		})

		// hardening of cookies, non-persistent cookies, only available for http and https in production builds
		cs := sessions.NewCookieStore([]byte(conf.SessionSecret))
		cs.Options.MaxAge = 0
		cs.Options.HttpOnly = true
		cs.Options.SameSite = http.SameSiteStrictMode
		cs.Options.Secure = config.IsProdBuild()
		app.SessionStore = cs

		// Log request parameters (filters apply).
		app.Use(paramlogger.ParameterLogger)

		// initialize our internal logging with buffalo logger
		if err := logger.SetLoggers(app.Logger); err != nil {
			log.Fatal(err)
		}

		// initialize the DB
		if err := models.ConnectDB(conf.Environment); err != nil {
			logger.BuffaloLogger.Fatal(err)
		}

		// Protect against CSRF attacks. https://www.owasp.org/index.php/Cross-Site_Request_Forgery_(CSRF)
		// Remove to disable this.
		app.Use(csrf.New)

		// Wraps each request in a transaction.
		//  c.Value("tx").(*pop.Connection)
		// Remove to disable this.
		app.Use(popmw.Transaction(models.DB))

		// URLs for user
		app.GET("/", AccountRequestPage).App.Redirect(302, "/accountrequest", "/")
		app.POST("/accountrequest", AccountRequestProcessRequest)
		app.GET("/emailvalidation/{token}", EmailValidationProcessRequest)
		app.GET("/credentialreset", CredentialsResetRequest)
		app.POST("/credentialreset/", CredentialsResetProcessRequest)
		app.GET("/credentialreset/process/{tokenpwr}", CredentialsResetProcessToken)
		app.POST("/credentialreset/process/{tokenpwr}/create", CredentialsResetProcessCredentials)

		// experimental things not in the production builds
		if !config.IsProdBuild() {
			app.GET("/user/ipa-account/validation/{ipa_token_account_validation}/", ValidateTokenIPA)
		}

		// URL for admin
		admin := app.Group("/admin")
		admin.Resource("/", RequestResource{})

		// custom assets if defined
		if customAssetsBox != nil {
			app.ServeFiles("/cassets", customAssetsBox)
		}

		app.ServeFiles("/", http.FS(public.FS())) // serve files from the public directory
	}

	return app
}
