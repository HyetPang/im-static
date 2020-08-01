package main

import (
	"im-static/controller"
	"os"
	"path/filepath"
	"time"

	"github.com/zengyu2020/mskit"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	if err := os.MkdirAll(controller.UploadFilePath, os.ModePerm); err != nil {
		zap.L().Fatal("检查或创建上传目录时发生错误，程序即将退出", zap.Error(err))
	}

	mskit.NewBuilder(
		mskit.UseRedis(),
		mskit.UseHttp(false, false),
		mskit.UseZap("logs", filepath.Base(os.Args[0]), mskit.MinZapLevel(zapcore.InfoLevel)),
		mskit.UseInvoke(registerHTTPHandlers),
	).Run()
}

func registerHTTPHandlers(router gin.IRouter, redisC redis.UniversalClient) {
	staticFileController := controller.NewStaticFileController()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AddAllowHeaders("X-Requested-With")
	corsConfig.AddAllowHeaders("Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
	corsConfig.AddExposeHeaders("Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar")
	corsConfig.MaxAge = 172800 * time.Second
	corsConfig.AllowCredentials = false
	corsConfig.AllowFiles = true
	router.Use(cors.New(corsConfig))

	router.GET("", staticFileController.Get(redisC))
	router.POST("/upload", staticFileController.Upload(redisC))
	//router.DELETE("/:id", staticFileController.Delete)
}
