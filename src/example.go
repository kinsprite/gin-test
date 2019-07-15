package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"go.elastic.co/apm/module/apmgin"
	"gopkg.in/resty.v1"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

const prefixV1 = "/api/gin/v1"
const prefixV2 = "/api/gin/v2"

var userServerURL = "http://user-test:80"

func init() {
	url := os.Getenv("USER_SERVER_URL")

	if url != "" {
		userServerURL = url
	}
}

func main() {
	engine := gin.New()
	engine.Use(apmgin.Middleware(engine))

	v1 := engine.Group(prefixV1)

	v1.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	v2 := engine.Group(prefixV2)

	v2.GET("/productsDetails", func(c *gin.Context) {
		resp, err := http.Get(userServerURL + "/api/user/v1/userInfoBySession")

		if err != nil {
			log.Println("ERROR   request user info")
			return
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			log.Println("ERROR   reading user info")
			return
		}

		var userInfo UserInfo
		json.Unmarshal(body, &userInfo)

		c.JSON(http.StatusOK, gin.H{
			"message":  "all products' details",
			"userId":   userInfo.ID,
			"userName": userInfo.Name,
		})
	})

	v2.GET("/userInfo", func(c *gin.Context) {
		resp, err := resty.R().Get(userServerURL + "/api/user/v1/userInfoBySession")

		if err != nil {
			log.Println("ERROR   request user info")
			return
		}

		var userInfo UserInfo
		json.Unmarshal(resp.Body(), &userInfo)

		c.JSON(http.StatusOK, gin.H{
			"userId":   userInfo.ID,
			"userName": userInfo.Name,
		})
	})

	engine.Run(":8080") // 监听并在 0.0.0.0:8080 上启动服务
}
