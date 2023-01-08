package models

import "time"

type User_File struct {
	Id          int `gorm:"primarykey;autoIncrement"`
	UserId      int
	File_sha1   string
	FileSize    int64
	FileName    string
	UploadTime  time.Time
	Last_Update time.Time `gorm:"autoUpdateTime"`
}
