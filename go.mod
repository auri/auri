module auri

go 1.16

require (
	github.com/gobuffalo/buffalo v0.16.16
	github.com/gobuffalo/buffalo-pop/v2 v2.3.0
	github.com/gobuffalo/envy v1.9.0
	github.com/gobuffalo/logger v1.0.3
	github.com/gobuffalo/mw-csrf v1.0.0
	github.com/gobuffalo/mw-paramlogger v1.0.0
	github.com/gobuffalo/nulls v0.4.0
	github.com/gobuffalo/packr/v2 v2.8.1
	github.com/gobuffalo/pop/v5 v5.3.1
	github.com/gobuffalo/suite/v3 v3.0.0
	github.com/gobuffalo/validate/v3 v3.3.0
	github.com/gobuffalo/x v0.1.0
	github.com/gofrs/uuid v4.0.0+incompatible
	github.com/gorilla/sessions v1.2.1
	github.com/markbates/grift v1.5.0
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.7.0
	github.com/stanislas-m/mocksmtp v1.1.0
	github.com/stretchr/testify v1.6.1
	github.com/tehwalris/go-freeipa v0.0.0-20200322083409-e462fc554b76
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad
)

replace github.com/tehwalris/go-freeipa => ./support/go-freeipa
