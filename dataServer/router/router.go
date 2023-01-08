package router

import (
	"file-server/dataServer/controller/file"
	"file-server/dataServer/controller/tempFile"
	"file-server/token"

	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine) {
	fileGroup := r.Group("/file", token.AuthMiddleware())
	{
		fileGroup.PUT("/upload", file.Upload)
		fileGroup.GET("/download", file.Download)

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

	tempFileGroup := r.Group("/temp")
	{
		tempFileGroup.POST("/fileMeta", tempFile.UploadMeta)        //上传元数据
		tempFileGroup.PATCH("/file", tempFile.UploadtoTempFile)     //上传到暂时文件
		tempFileGroup.PUT("/removeToStore", tempFile.RemoveToStore) //转正
		tempFileGroup.DELETE("/fileDelete", tempFile.DeleteFile)
	}
}
