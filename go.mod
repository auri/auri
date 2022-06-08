module auri

go 1.16

require (
	github.com/gobuffalo/buffalo v0.18.4
	github.com/gobuffalo/buffalo-pop/v3 v3.0.3
	github.com/gobuffalo/envy v1.10.1
	github.com/gobuffalo/flect v0.2.5 // indirect
	github.com/gobuffalo/logger v1.0.6
	github.com/gobuffalo/mw-csrf v1.0.0
	github.com/gobuffalo/mw-paramlogger v1.0.0
	github.com/gobuffalo/nulls v0.4.1
	github.com/gobuffalo/packr/v2 v2.8.1
	github.com/gobuffalo/pop/v5 v5.3.4 // indirect
	github.com/gobuffalo/pop/v6 v6.0.2
	github.com/gobuffalo/suite/v3 v3.0.2
	github.com/gobuffalo/suite/v4 v4.0.2
	github.com/gobuffalo/validate/v3 v3.3.1
	github.com/gobuffalo/x v0.0.0-20190224155809-6bb134105960
	github.com/gofrs/uuid v4.2.0+incompatible
	github.com/gorilla/sessions v1.2.1
	github.com/markbates/grift v1.5.0
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.8.1
	github.com/stanislas-m/mocksmtp v1.1.0
	github.com/stretchr/testify v1.7.1
	github.com/tehwalris/go-freeipa v0.0.0-00010101000000-000000000000
	golang.org/x/crypto v0.0.0-20211117183948-ae814b36b871
)

replace github.com/tehwalris/go-freeipa => ./support/go-freeipa
