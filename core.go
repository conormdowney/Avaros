package main

import (
	"context"
	"fmt"
	"net/http"

	rest "avaros/rest"
	router "avaros/router"

	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	fmt.Println("server running")

	router := router.NewRouter()
	db := newDatabase()
	defer db.Close()

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

	http.ListenAndServe("localhost:3000", router)
}

// newDatabase connects to a database at the start and passes that connection to
// any service below it.
// Ideally this would be a wrapper that uses an interface. This way you coud have different implementations
// for different database types if needed. In the interest of time ive skipeed this for now
// use a Pool as it is thread safe whereas the standard connection is not
func newDatabase() *pgxpool.Pool {
	// database connection string is in db_config.txt
	// Would have used environemnt variables but i dont know if there are issues with adding them
	// or if whoever runs this would just prefer not to have them on their machine,
	// so I just used a text file instead
	// dbConfig, err := os.ReadFile("db_config.txt")
	// if err != nil {
	// 	panic(err.Error())
	// }
	//conn, err := pgx.NewConnPool(context.Background(), string(dbConfig))
	conn, err := pgxpool.Connect(context.Background(), "user=postgres dbname=postgres password=password host=localhost sslmode=disable")
	if err != nil {
		panic("Unable to connect to database: " + err.Error())
	}

	return conn
}
