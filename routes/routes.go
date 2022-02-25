package routes

import (
	"etl-blocks2/controllers"

	"github.com/gin-gonic/gin"
)

func LoadRoutes() {

	router := gin.Default()

	router.GET("/", controllers.GetTest)

	router.GET("/data", controllers.GetData)
	
	router.NoRoute(noRoute)

	router.Run()

}

func noRoute (c *gin.Context) {
	c.JSON(404, gin.H{
		"message": "Page not found."
	})
}
