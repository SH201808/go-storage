package mysql

import (
	"file-server/models"
	"file-server/models/meta"
	"file-server/setting"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init(cfg *setting.MySQLConfig) (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Root,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("gorm Open failed,err:", err)
		return
	}
	DB = db
	sqlDB, err := DB.DB()
	if err != nil {
		log.Println("gorm DB failed,err:", err)
		return
	}
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConn)
	sqlDB.SetMaxOpenConns(cfg.MaxIdleConn)

	DB.AutoMigrate(&meta.File{}, models.User{}, models.User_File{})

	return
}
