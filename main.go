package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	database "avaros/database"
	rest "avaros/rest"
	router "avaros/router"

	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	fmt.Println("server running")

	router := router.NewRouter()
	db := newDatabase()
	defer db.Close()
	// create and seed the database
	database.Seed(db)

	// instantiate a rest object so all rest services have the same database and router
	RestObj := rest.RestServiceObject{
		router,
		db,
	}

	// Doing it this way as it is easy then to add any more services as required
	restServices := []rest.RestService{
		&rest.RoomService{RestObj},
	}

	// Loop through and initialise their routes
	for _, service := range restServices {
		err := service.Init()
		if err != nil {
			panic(err)
		}
	}

	listen := os.Getenv("LISTEN_ADDR")
	fmt.Println("LISTEN: " + listen)
	if listen == "" {
		listen = "localhost:3000"
	}

	http.ListenAndServe(listen, router)
}

// newDatabase connects to a database at the start and passes that connection to
// any service below it.
func newDatabase() *pgxpool.Pool {

	user := os.Getenv("DB_USER")
	if user == "" {
		user = "postgres"
	}

	name := os.Getenv("DB_NAME")
	if name == "" {
		name = "postgres"
	}

	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "password"
	}

	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}

	conn, err := pgxpool.Connect(context.Background(),
		fmt.Sprintf("user=%s dbname=%s password=%s host=%s sslmode=disable",
			user, name, password, host))
	if err != nil {
		panic("Unable to connect to database: " + err.Error())
	}

	return conn
}
