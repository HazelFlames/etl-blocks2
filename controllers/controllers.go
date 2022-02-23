package controllers

import (
	"etl-blocks2/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetData(c *gin.Context) {
	areas := models.ReadPg()
	c.IndentedJSON(http.StatusOK, areas)
}

func GetTest(c *gin.Context) {

	c.IndentedJSON(http.StatusOK, "")
}
