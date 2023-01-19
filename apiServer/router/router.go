package router

import (
	file "file-server/apiServer/controller/File"
	"file-server/apiServer/controller/Resume"
	user "file-server/apiServer/controller/User"
	"file-server/token"

	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine) {
	fileGroup := r.Group("/file", token.AuthMiddleware())
	{
		fileGroup.PUT("/upload", file.Upload)
		fileGroup.GET("/download", file.Download)

		resumeGroup := fileGroup.Group("/resume")
		{
			resumeGroup.POST("/upload", Resume.Upload)
		}

		// fileGroup.GET("/query", File.Query)
		// fileGroup.DELETE("/delete", File.Delete)
		// fileGroup.PUT("/update", File.Update)
		// fileGroup.POST("/tryFastUpload", File.TryFastUpload)

		// mpUpload := fileGroup.Group("/mpupload")
		// {
		// 	mpUpload.POST("/init", fileblock.InitMeta)
		// 	mpUpload.POST("/uppart", fileblock.Uppart)
		// 	mpUpload.POST("/complete", fileblock.CompleteUpload)
		// }
	}

	userGroup := r.Group("/user")
	{
		userGroup.POST("/register", user.Register)
		userGroup.POST("/login", user.Login)
		// userGroup.GET("/query", User.Query)
		// userGroup.GET("/getAccessToken", User.GetAccessToken)
	}
}
