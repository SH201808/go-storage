package main

import (
	"file-server/dataServer/heartBeat"
	"file-server/dataServer/locate"
	"file-server/dataServer/router"
	"file-server/setting"
	"log"

	"github.com/gin-gonic/gin"
)

func init() {
	log.SetFlags(log.Llongfile)

	setting.Init()
}

func main() {
	locate.SetFileLoc()
	go heartBeat.Start()
	go locate.Start()
	r := gin.Default()

	router.Setup(r)
	r.Run(setting.Conf.Port)
}
