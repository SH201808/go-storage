package router

import (
	"file-server/controller/File"
	fileblock "file-server/controller/FileBlock"
	"file-server/controller/User"
	"file-server/token"

	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine) {

	fileGroup := r.Group("/file", token.AuthMiddleware())
	{
		fileGroup.POST("/upload", File.Upload)
		fileGroup.GET("/query", File.Query)
		fileGroup.GET("/download", File.Download)
		fileGroup.DELETE("/delete", File.Delete)
		fileGroup.PUT("/update", File.Update)
		fileGroup.POST("/tryFastUpload", File.TryFastUpload)

		mpUpload := fileGroup.Group("/mpupload")
		{
			mpUpload.POST("/init", fileblock.InitMeta)
			mpUpload.POST("/uppart", fileblock.Uppart)
			mpUpload.POST("/complete", fileblock.CompleteUpload)
		}
	}

	userGroup := r.Group("/user")
	{
		userGroup.POST("/register", User.Register)
		userGroup.POST("/login", User.Login)
		userGroup.GET("/query", User.Query)
		userGroup.GET("/getAccessToken", User.GetAccessToken)
	}
}
