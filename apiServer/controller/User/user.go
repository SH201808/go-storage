package user

import (
	"file-server/dao"
	response "file-server/models/Response"
	"file-server/token"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	userName := c.PostForm("userName")
	password := c.PostForm("password")

	err := dao.AddUser(userName, password)
	if err != nil {
		c.JSON(http.StatusOK, response.Err("注册失败"))
	}
	c.JSON(http.StatusOK, response.Success("注册成功"))
}

func Login(c *gin.Context) {
	userName := c.PostForm("userName")
	password := c.PostForm("password")

	user := dao.QueryUser(userName, password)
	if user == nil {
		c.JSON(http.StatusOK, response.Success("未找到用户"))
		return
	}

	AccessToken, err := token.AccessCreate(user.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Err("生成Access token 错误"))
		log.Println("生成Access token err：", err)
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"AccessToken": AccessToken,
		"expireAt":    time.Now().Add(token.AccessTokenExpiredTime),
	}))
}
