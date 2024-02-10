package main

import (
	"MedicalLowCode-backend/router"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	//静态文件服务
	r.Static("/static", "./static")

	router.ApiRouterInit(r)
	r.Run(":8080") // 监听并在 0.0.0.0:8000 上启动服务
}
