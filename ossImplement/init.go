package ossImplement

import (
	"file-server/setting"
	"fmt"
	"os"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type FileMeta struct {
	AbsLoc   string
	FileSize int64
}

func ConstructFileMeta(absLoc string, fileSize int64) FileMeta {
	return FileMeta{
		AbsLoc:   absLoc,
		FileSize: fileSize,
	}
}

var Client *oss.Client

func Init(config *setting.OssConfig) {
	client, err := oss.New(config.EndPoint, config.EndPoint, config.AccessKeySecret)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}
	Client = client
}
