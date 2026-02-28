package router

import (
	"github.com/gin-gonic/gin"
	"go-file/controller"
	"go-file/middleware"
)

func setApiRouter(router *gin.Engine) {
	router.Use(middleware.GlobalAPIRateLimit())
	router.GET("/status", controller.GetStatus)
	router.POST("/api/file", middleware.FileUploadPermissionCheck(), controller.UploadFile)
	router.POST("/api/image", middleware.ImageUploadPermissionCheck(), controller.UploadImage)
	router.GET("/api/notice", controller.GetNotice)

	// Chat API (requires login)
	chatAuth := router.Group("/api/chat")
	chatAuth.Use(middleware.ApiAuth())
	{
		chatAuth.GET("/messages", controller.GetChatMessages)
		chatAuth.POST("/send", controller.SendChatMessage)
		chatAuth.POST("/upload", controller.UploadChatFile)
		chatAuth.GET("/poll", controller.PollChatMessages)
	}

	basicAuth := router.Group("/api")
	basicAuth.Use(middleware.ApiAuth())
	{
		basicAuth.DELETE("/file", controller.DeleteFile)
		basicAuth.DELETE("/image", controller.DeleteImage)
		basicAuth.PUT("/user", middleware.NoTokenAuth(), controller.UpdateSelf)
		basicAuth.POST("/token", controller.GenerateNewUserToken)
	}
	adminAuth := router.Group("/api")
	adminAuth.Use(middleware.ApiAdminAuth())
	{
		adminAuth.POST("/user", controller.CreateUser)
		adminAuth.PUT("/manage_user", controller.ManageUser)
		adminAuth.GET("/option", controller.GetOptions)
		adminAuth.PUT("/option", controller.UpdateOption)
		adminAuth.DELETE("/chat/clear", controller.ClearChatMessages)
		statRouter := adminAuth.Group("/stat")
		{
			statRouter.GET("/ip", controller.GetIPs)
			statRouter.GET("/url", controller.GetURLs)
			statRouter.GET("/req", controller.GetReqs)
		}
	}
}