package main

import (
	heartbeat "file-server/apiServer/heartBeat"
	"file-server/apiServer/router"
	"file-server/db/mysql"
	"file-server/elasticSearch"
	"file-server/models/meta"
	"file-server/setting"
	"log"

	"github.com/gin-gonic/gin"
)

func init() {
	log.SetFlags(log.Llongfile)

	setting.Init()
	mysql.Init(setting.Conf.MySQLConfig)

	elasticSearch.Init(setting.Conf.ElasticSearchConfig)
	err := elasticSearch.CreateIndex("filemeta", meta.Mappings)
	if err != nil {
		log.Println("create index err: ", err)
		return
	}
}

func main() {
	go heartbeat.Listen()

	r := gin.Default()
	router.Setup(r)

	r.Run(":9999")
}
