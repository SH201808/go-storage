package main

import (
	"file-server/db/mysql"
	"file-server/db/redis"
	"file-server/ossImplement"
	rabbitmq "file-server/rabbitMQ"
	"file-server/router"
	"file-server/setting"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	//create("../cmd/test.txt")
	// block("../cmd/test.txt")
	setting.Init()

	mysql.Init(setting.Conf.MySQLConfig)
	redis.Init(setting.Conf.RedisConfig)
	ossImplement.Init(setting.Conf.OssConfig)
	rabbitmq.Conn = rabbitmq.Init(setting.Conf.RabbitMQConfig)
	defer CloseAllConn()

	r := gin.Default()
	router.Setup(r)

	r.Run()
}

func CloseAllConn() {
	redis.DB.Close()
	rabbitmq.Conn.Close()
}

func create(dst string) {
	file, err := os.Create(dst)
	if err != nil {
		log.Println("创建测试文件错误")
		return
	}
	defer file.Close()
	var data strings.Builder

	for i := 0; i < 1024*1024; i++ {
		data.WriteString("hhhhhhhhhhh")
	}
	_, err = file.WriteString(data.String())
	if err != nil {
		return
	}
}

func block(dst string) {
	file, err := os.Open(dst)
	if err != nil {
		log.Println("测试打开文件错误")
		return
	}
	defer file.Close()
	fileInfo, _ := file.Stat()
	fmt.Println(fileInfo.Size())
	chunckCount := int(math.Ceil(float64(fileInfo.Size()) / (5 * 1024 * 1024)))
	data, err := ioutil.ReadFile(dst)
	if err != nil {
		log.Println("读取测试文件错误")
		return
	}
	for i := 0; i < chunckCount; i++ {
		seq := strconv.Itoa(i)
		newFile, err := os.Create("../cmd/" + seq)
		if err != nil {
			log.Println("创建分块文件错误")
			return
		}
		defer newFile.Close()
		if i == chunckCount-1 {
			newFile.Write(data[i*5*1024*1024:])
		} else {
			newFile.Write(data[i*5*1024*1024 : (i+1)*5*1024*1024])
		}
	}
}
