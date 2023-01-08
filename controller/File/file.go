package File

import (
	"encoding/json"
	"file-server/dao"
	response "file-server/models/Response"
	"file-server/models/meta"
	"file-server/ossImplement"
	rabbitmq "file-server/rabbitMQ"
	"file-server/setting"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		log.Println("upload file error:", err.Error())
		c.JSON(http.StatusInternalServerError, response.Err("upload file error"))
		return
	}

	dst := setting.Conf.AbsLoc + file.Filename
	//判断本地存储是否存在文件
	if _, err := os.Stat(dst); err != nil {
		err = c.SaveUploadedFile(file, dst)
		if err != nil {
			log.Println("Save file error:", err.Error())
			c.JSON(http.StatusInternalServerError, response.Err("Save file to local error"))
			return
		}
	}
	//数据库中是否已经存在文件
	userId, _ := c.Get("userId")
	fileMeta := *meta.ConstructeFile(dst)
	err = dao.UploadFileMetaAndUserFile(fileMeta, userId.(int))

	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Err("Save file to database error"))
	}

	c.JSON(http.StatusOK, response.Success())

	q := rabbitmq.ConstructQueue("toOss", false, false, false, false, nil)

	ossFileMeta := ossImplement.ConstructFileMeta(dst, file.Size)
	data, err := json.Marshal(ossFileMeta)
	publish := rabbitmq.ConstructPublish("", q.Name, false, false, data)

	rabbitmq.SendMsg(rabbitmq.Conn, q, publish)
}

func Download(c *gin.Context) {
	fileSha := c.Query("fileSha")
	file := dao.GetFile(fileSha)
	if file.FileName == "" {
		log.Println("File is not exist")
		c.JSON(http.StatusInternalServerError, response.Err("File is not exist"))
		return
	}
	dst := "../fileStore/" + file.FileName
	c.File(dst)
	//c.JSON(http.StatusOK, response.Success())
}

func Delete(c *gin.Context) {
	fileSha := c.Query("fileSha")
	file := dao.GetFile(fileSha)
	if file.FileName == "" {
		log.Println("File is not exist")
		c.JSON(http.StatusInternalServerError, response.Err("File is not exist"))
		return
	}
	dst := "../fileStore/" + file.FileName

	os.Remove(dst)
	dao.DeleteFile(fileSha)
	c.JSON(http.StatusOK, response.Success())
}

func Update(c *gin.Context) {
	fileSha := c.Query("fileSha")
	newName := c.Query("newName")
	file := dao.GetFile(fileSha)
	if file.FileName == "" {
		log.Println("File is not exist")
		c.JSON(http.StatusInternalServerError, response.Err("File is not exist"))
		return
	}
	dst := "../fileStore/" + file.FileName
	newDst := "../fileStore/" + newName
	file.FileName = newDst
	userId, _ := c.Get("userId")
	dao.UploadFileMetaAndUserFile(*file, userId.(int))

	os.Rename(dst, newDst)
	c.JSON(http.StatusOK, response.Success())
}

func Query(c *gin.Context) {
	userId, _ := c.Get("userId")
	userFileMetas := dao.QueryUserFileMetas(userId.(int))
	if len(userFileMetas) == 0 {
		c.JSON(http.StatusOK, response.Err("没有文件"))
		return
	}
	c.JSON(http.StatusOK, response.Success(gin.H{
		"data": gin.H{
			"Files": userFileMetas,
		},
	}))
}

func TryFastUpload(c *gin.Context) {
	fileSha1 := c.PostForm("sha1")

	//从文件表中查询相同hash的文件记录
	fileMeta := dao.GetFile(fileSha1)

	if fileMeta.Dst == "" {
		c.JSON(http.StatusOK, response.Success("秒传失败，请选择正常上传接口"))
		return
	}

	userId := c.GetInt("userId")
	err := dao.UploadFileMetaAndUserFile(*fileMeta, userId)
	if err != nil {
		c.JSON(http.StatusOK, response.Err("上传用户文件信息错误"))
		return
	}

	c.JSON(http.StatusOK, response.Success("秒传成功"))
}
