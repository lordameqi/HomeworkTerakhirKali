package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
)

type User struct {
	Username string `json:"username"`
	Follower int    `json:"followers"`
}

func getJson() []byte {

	json_resp, err := http.Get("https://jsonkeeper.com/b/DMXK")
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	json_data, err := ioutil.ReadAll(json_resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return json_data

}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	router.GET("/:userid", func(c *gin.Context) {
		keyword := c.Param("userid")

		user := map[string]User{}
		json.Unmarshal(getJson(), &user)

		c.JSON(http.StatusOK, user[keyword])
	})

	router.GET("/followers/:username", func(c *gin.Context) {
		keyword := c.Param("username")

		user := map[string]User{}
		json.Unmarshal(getJson(), &user)

		result := 0
		for _, val := range user {
			if val.Username == keyword {
				result = val.Follower
			}
		}

		c.JSON(http.StatusOK, result)
	})

	router.Run(":" + port)
}
