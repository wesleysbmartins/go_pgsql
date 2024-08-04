package main

import (
	"go_pgsql/internal/adapters/database/sql/postgres"
)

func init() {
	postgres := postgres.Postgres{}
	postgres.Connect()
}

func main() {

}
