package handler

import (
	"github.com/ReisenCW/go-url-shortener/shortener"
	"github.com/ReisenCW/go-url-shortener/store"
	"github.com/gin-gonic/gin"
	"net/http"
)


/*
暂存从客户端请求中解析出的数据: 长url与用户ID
c.ShouldBindJSON(&creationRequest) ：将请求体中的JSON数据绑定到结构体
	将JSON字段映射到struct中（如 JSON 的 long_url → struct的 LongUrl）
	binding:"required" ：检查 该字段 是否存在，若缺失则返回错误
*/
type UrlCreationRequest struct {
	LongUrl string `json:"long_url" binding:"required"`
	UserId  string `json:"user_id" binding:"required"`
}

func CreateShortUrl(c *gin.Context) {
	var creationRequest UrlCreationRequest

	err := c.ShouldBindJSON(&creationRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shortUrl := shortener.GenerateShortLink(creationRequest.LongUrl, creationRequest.UserId)
	store.SaveUrlMapping(shortUrl, creationRequest.LongUrl, creationRequest.UserId)

	// 通过 c.JSON 返回 200 状态码（HTTP 标准的 “成功” 状态）
	// 并以 JSON 格式返回响应：包含message和完整的short_url
	host := "http://localhost:9808/"
	c.JSON(200, gin.H{
		"message":   "short url created successfully",
		"short_url": host + shortUrl,
	})
}

func HandleShortUrlRedirect(c *gin.Context) {
	// c.Param: 获取 URL 路径中定义的 “路径参数”
	// main中的路由定义 r.GET("/:shortUrl", ...)，:
	// shortUrl 是一个动态路径参数（类似占位符），代表短url的唯一标识
	shortUrl := c.Param("shortUrl")
	initialUrl := store.RetrieveInitialUrl(shortUrl)
	// 302 HTTP 状态码，表示 "临时重定向"
	c.Redirect(302, initialUrl)
}