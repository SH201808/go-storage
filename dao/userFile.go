package dao

import (
	"file-server/db/mysql"
	"file-server/models"
	"log"
)

func QueryUserFileMetas(userId int) (userFileMetas []models.User_File) {
	err := mysql.DB.Where("user_id = ?", userId).Find(&userFileMetas).Error
	if err != nil {
		log.Println("QueryUserFileMeats err:", err)
		return
	}
	return
}
