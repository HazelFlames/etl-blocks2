package main

import (
	"etl-blocks2/db"
	"etl-blocks2/models"
)

func main() {
	db.DbConect()
	models.ReadPg()

}
