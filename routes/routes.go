package routes

import (
	"etl-blocks2/controllers"

	"github.com/gin-gonic/gin"
)

func LoadRoutes() *gin.Engine {

	router := gin.Default()

	router.GET("/data", controllers.GetData)

	router.GET("", controllers.GetTest)

	return router
}
