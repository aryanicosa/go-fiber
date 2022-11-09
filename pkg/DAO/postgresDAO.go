package DAO

import (
	"github.com/aryanicosa/go-fiber-rest-api/app/queries"
)

type DB interface {
	InitDBConnection() error
	UserConn() (*queries.UserQueries, error)
	BookConn() (*queries.BookQueries, error)
}
