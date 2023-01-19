package Resume

import (
	heartbeat "file-server/apiServer/heartBeat"
	"file-server/apiServer/locate"
	response "file-server/models/Response"
	"file-server/rs"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

func Upload(c *gin.Context) {
	hash := c.Request.Header.Get("Digests")
	if locate.Exist(hash) {
		// todo 新增文件版本号
		c.JSON(http.StatusOK, response.Success("save file success"))
		return
	}
	// c.Request.URL.EscapedPath()
	// todo 上传携带空格的数据存在转义问题
	name := c.Query("fileName")
	FileSize := c.Request.Header.Get("FileSize")
	fileSize, _ := strconv.Atoi(FileSize)
	storeResumeObject(c, name, hash, int64(fileSize))
	return
}

func storeResumeObject(c *gin.Context, name string, hash string, size int64) {
	dataServers := heartbeat.ChooseDataServers(rs.ALL_SHARDS, nil)
	if len(dataServers) != rs.ALL_SHARDS {
		log.Println("can not find enough dataServer")
		c.Status(http.StatusServiceUnavailable)
		return
	}
	stream, err := rs.NewResumablePutStream(dataServers, name, hash, size)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Header("location", "/resume/"+url.PathEscape(stream.ToToken()))
	c.Status(http.StatusCreated)
}
