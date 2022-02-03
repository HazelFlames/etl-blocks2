package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func DbConect() *sql.DB {

	password := os.Getenv("PG_ETL_PASS")
	user := os.Getenv("PG_ETL_USER")
	dbName := os.Getenv("PG_ETL_DB")
	host := os.Getenv("PG_ETL_HOST")

	connection := fmt.Sprintf("user=%s dbname=%s password=%s host=%s sslmode=disable", user, dbName, password, host)

	db, err := sql.Open("postgres", connection)
	if err != nil {
		log.Println("Unable to connect:" + err.Error())
	} else {
		log.Println("Postgres database connected!")
	}
	return db
}
