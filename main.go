package main

import (
	"fmt"
	"log"
	"net/http"
	"open/service"
	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("start")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	router := gin.Default()
	router.POST("/", func(c *gin.Context) {
		fmt.Println("POST", "/")
		c.String(http.StatusOK, "Hello World test2.go,/")
	})
	router.GET("/", func(c *gin.Context) {
		fmt.Println("POST", "/")
		c.String(http.StatusOK, "Hello World test2.go,/")
	})
	router.POST("/test",  service.ServiceAll )
	router.GET("/test", func(c *gin.Context) {
		fmt.Println("POST", "test")
		c.String(http.StatusOK, "Hello World test2.go,/test")
	})
	router.GET("/oath/redirect", func(c *gin.Context) {
		fmt.Println("GET", "/oath/redirect")
		c.String(http.StatusOK, "Hello World test2.go,/test")
	})
	router.Run(":8000")
}
