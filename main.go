package main

import (
	"MedicalLowCode-backend/router"
	"MedicalLowCode-backend/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	r := gin.Default()

	// 静态文件服务
	//r.Static("/static", "./static")

	// 跨域请求处理
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		//c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}
		c.Next()
	})

	//文件鉴权
	r.GET("/fileAuth", router.CheckToken, func(c *gin.Context) {
		token := c.MustGet("token").(*util.Token)
		filePath := c.Query("path")
		pathParts := strings.Split(filePath, "/")
		var userId string
		if len(pathParts) >= 3 {
			userId = pathParts[2]
		} else {
			c.JSON(http.StatusForbidden, gin.H{})
			return
		}
		if userId != strconv.Itoa(token.UserId) {
			c.JSON(http.StatusUnauthorized, gin.H{})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{})
			return
		}
	})

	router.ApiRouterInit(r)
	r.Run(":8080") // 监听并在 0.0.0.0:8000 上启动服务
}
