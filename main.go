package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"score-summary-backend/routers"
)

func main() {

	gin.SetMode(gin.ReleaseMode)
	// 创建 Gin 实例
	r := gin.Default()

	// 加载路由
	routers.LoadUsersRoutes(r)

	// 启动服务器
	log.Println("start")
	if err := r.Run(":8888"); err != nil {
		panic("Failed to start server!")
	}
}
