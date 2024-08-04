package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Postgres struct {
	host     string
	port     int
	user     string
	password string
	database string
}

type IPostgres interface {
	Connect()
}

var Client *sql.DB

func (p *Postgres) Connect() {

	credentials := Postgres{
		host:     "localhost",
		port:     5432,
		user:     "postgres",
		password: "1234",
		database: "teste",
	}

	connStr := fmt.Sprintf("host=%s port=%v user=%s password=%s dbname=%s sslmode=disable", credentials.host, credentials.port, credentials.user, credentials.password, credentials.database)

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		panic(err)
	}

	err = db.Ping()

	if err != nil {
		panic(err)
	}

	fmt.Println("POSTGRES CONNECTION SUCCESS!")

	Client = db
}
