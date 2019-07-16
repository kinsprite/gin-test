package main

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmgin"
	"go.elastic.co/apm/module/apmgrpc"
	"go.elastic.co/apm/module/apmhttp"
	"golang.org/x/net/context/ctxhttp"
	"google.golang.org/grpc"

	// "gopkg.in/resty.v1"
	pb "github.com/kinsprite/producttest/pb"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary
var tracingClient = apmhttp.WrapClient(http.DefaultClient)

const prefixV1 = "/api/gin/v1"
const prefixV2 = "/api/gin/v2"

var userServerURL = "http://user-test:80"
var productServerAddress = "producttest:80"

const (
	defaultName = "world"
)

func init() {
	url := os.Getenv("USER_SERVER_URL")

	if url != "" {
		userServerURL = url
	}

	address := os.Getenv("PRODUCT_SERVER_ADDRESS")

	if address != "" {
		productServerAddress = address
	}
}

func fetchProduct() string {
	// Set up a connection to the server.
	conn, err := grpc.Dial(productServerAddress, grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(apmgrpc.NewUnaryClientInterceptor()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.Message)

	result := r.Message

	r, err = c.SayHelloAgain(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.Message)
	return result
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
		req := c.Request
		resp, err := ctxhttp.Get(req.Context(), tracingClient, userServerURL+"/api/user/v1/userInfoBySession")

		if err != nil {
			apm.CaptureError(req.Context(), err).Send()
			log.Println("ERROR   request user info")
			c.AbortWithError(500, errors.New("failed to query backend"))
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

		productMsg := fetchProduct()

		c.JSON(http.StatusOK, gin.H{
			"message":    "all products' details",
			"userId":     userInfo.ID,
			"userName":   userInfo.Name,
			"productMsg": productMsg,
		})
	})

	// v2.GET("/userInfo", func(c *gin.Context) {
	// 	resp, err := resty.R().Get(userServerURL + "/api/user/v1/userInfoBySession")

	// 	if err != nil {
	// 		log.Println("ERROR   request user info")
	// 		return
	// 	}

	// 	var userInfo UserInfo
	// 	json.Unmarshal(resp.Body(), &userInfo)

	// 	c.JSON(http.StatusOK, gin.H{
	// 		"userId":   userInfo.ID,
	// 		"userName": userInfo.Name,
	// 	})
	// })

	engine.Run(":8080") // 监听并在 0.0.0.0:8080 上启动服务
}
