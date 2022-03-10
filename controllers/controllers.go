package controllers

import (
	"encoding/json"
	"etl-blocks2/dbRedis"
	"etl-blocks2/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func PostData(c *gin.Context) {
	areas := models.ReadPg()
	c.IndentedJSON(http.StatusOK, areas)
}

func GetTest(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "up",
	})
}

func GetSize(c *gin.Context) {
	redis := dbRedis.ConnectRedis()
	defer redis.Close()

	dbsize, err := redis.DbSize().Result()
	if err != nil {
		log.Println(err.Error())
	}
	c.JSON(200, gin.H{
		"size": dbsize,
	})

}

func GetValues(c *gin.Context) {
	redis := dbRedis.ConnectRedis()
	defer redis.Close()

	keys, _ := redis.Keys("*").Result()

	var values []models.BQData
	var blocks models.BQData
	for _, v := range keys {
		value, err := redis.Get(v).Result()
		if err != nil {
			log.Println(err.Error())
		}
		json.Unmarshal([]byte(value), &blocks)
		values = append(values, blocks)
	}

	c.IndentedJSON(http.StatusOK, values)
}
