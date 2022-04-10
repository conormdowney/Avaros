/*
	Interface for rest services.
*/
package rest

import (
	"github.com/gocraft/web"
	"github.com/jackc/pgx/v4/pgxpool"
)

// RestService contains all the functions a rest service needs
type RestService interface {
	Init() error
}

// RestServiceObject contains the objects a rest service needs to run
type RestServiceObject struct {
	Router *web.Router
	Db     *pgxpool.Pool
}
