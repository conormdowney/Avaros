/*
	Wrapper class for the routing
*/

package router

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gocraft/web"
)

type Context struct {
	UserId int32
	Token  string
}

type User struct {
	Id int `json:"userId"`
}

// NewRouter creates a new instance of a web router for handling routing of paths
func NewRouter() *web.Router {
	router := web.New(Context{})
	router.Error((*Context).Error)
	router.Middleware((*Context).UserAuthentication)
	return router
}

// UserAuthentication is a (very) dummy function that would authenticate the user.
// I use userid instead of token for the sake of ease for someone running the code
func (c *Context) UserAuthentication(rw web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	cookie, err := r.Cookie("userId")
	if err != nil {
		panic("Error getting user id cookie")
	}
	userId, err := strconv.ParseInt(cookie.Value, 10, 32)
	if err != nil {
		panic("Error parsing user id")
	}

	userConfirmed := authenticateUser(userId)
	if userConfirmed {
		c.UserId = int32(userId)
		next(rw, r)
	} else {
		rw.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(rw, `{"error": "You are not authorized to use this service"}`)
	}
}

// Error handles any panics from within a rest call and lets you set an error code
// and the error message returned. Probably want to send back a user friendly error message but i
// just left it as the error as returned for this
func (c *Context) Error(rw web.ResponseWriter, r *web.Request, err interface{}) {
	rw.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(rw, `{"error": "`+err.(string)+`"}`)
}

func authenticateUser(userId int64) bool {
	if userId == 1 {
		return true
	} else {
		return false
	}
}
