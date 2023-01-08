package dao

import (
	"file-server/db/mysql"
	"file-server/models"
	"file-server/utils"
)

func AddUser(userName, Password string) error {
	user := &models.User{
		Name:     userName,
		Password: utils.GenSha1([]byte(Password)),
	}
	return mysql.DB.Create(user).Error
}

func QueryUser(userName, Password string) *models.User {
	user := models.User{}
	encodingPwd := utils.GenSha1([]byte(Password))
	mysql.DB.Where("name = ? AND password = ?", userName, encodingPwd).First(&user)
	return &user
}
