package main

import (
	"etl-blocks2/routes"
)

func main() {

	r := routes.LoadRoutes()

	r.Run("localhost:8080")

}
