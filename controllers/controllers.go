package controllers

import (
	"encoding/json"
	"etl-blocks2/dbRedis"
	"etl-blocks2/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	geojson "github.com/paulmach/go.geojson"
)

type Blocks struct {
	Block_id     int              `json:"block_id"`
	ClientBQ_id  int              `json:"client_id"`
	Block_parent int              `json:"block_parent"`
	Block_name   string           `json:"block_name"`
	Block_bounds geojson.Geometry `json:"bounds"`
	Block_abrv   string           `json:"abvr"`
}

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

	var values []Blocks
	var blocks Blocks
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
