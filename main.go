package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ReisenCW/go-url-shortener/handler"
	"github.com/ReisenCW/go-url-shortener/store"
)

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hey Go URL Shortener !",
		})
	})

	r.POST("/create-short-url", func(c *gin.Context) {
		handler.CreateShortUrl(c)
	})

	r.GET("/:shortUrl", func(c *gin.Context) {
		handler.HandleShortUrlRedirect(c)
	})

	// Note that store initialization happens here
	store.InitializeStore()


	err := r.Run(":9808")
	if err != nil {
		panic(fmt.Sprintf("Failed to start the web server - Error: %v", err))
	}
}

/*
使用命令行工具 curl 测试 POST 请求

curl -X POST http://localhost:9808/create-short-url \
   -H "Content-Type: application/json" \
   -d '{
   	"long_url": "https://www.eddywm.com/lets-build-a-url-shortener-in-go-part-iv-forwarding/",
   	"user_id": "ReisenCW"
   }'

得到返回的消息后，点击消息中的 short_url 链接，浏览器会跳转到 long_url
*/