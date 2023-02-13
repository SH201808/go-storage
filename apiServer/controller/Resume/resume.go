package Resume

import (
	heartbeat "file-server/apiServer/heartBeat"
	"file-server/apiServer/locate"
	"file-server/middleware"
	response "file-server/models/Response"
	"file-server/models/meta"
	"file-server/rs"
	"file-server/service/filemeta"
	"file-server/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

func GetToken(c *gin.Context) {
	hash := middleware.GetHashFromHeader(c.Request.Header)
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
	c.Header("location", url.PathEscape(stream.ToToken()))
	c.Status(http.StatusCreated)
}

func Upload(c *gin.Context) {
	// todo token的获取方式未确定
	token := c.Query("token")
	stream, err := rs.NewRSResumablePutStreamFromToken(token)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusForbidden)
		return
	}
	current := stream.CurrentSize()
	if current == -1 {
		log.Println("currentSize illegal")
		c.Status(http.StatusNotFound)
		return
	}
	offset := middleware.GetOffsetFromHeader(c.Request.Header)
	if current != offset {
		log.Println(current)
		log.Println(offset)
		c.Status(http.StatusRequestedRangeNotSatisfiable)
		return
	}

	bytes := make([]byte, rs.BLOCK_SIZE)
	for {
		n, err := io.ReadFull(c.Request.Body, bytes)
		if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
			log.Println(err)
			c.Status(http.StatusInternalServerError)
			return
		}
		current += int64(n)
		if current > stream.Size {
			stream.Commit(false)
			log.Println("resumable put exceed size")
			c.Status(http.StatusForbidden)
			return
		}
		if n != rs.BLOCK_SIZE && current != stream.Size {
			return
		}
		stream.Write(bytes[:n])
		if current == stream.Size {
			stream.Flush()
			getStream, err := rs.NewRSResumableGetStream(stream.Servers, stream.UUIDS, stream.Size)
			if err != nil {
				log.Println("new resumableGetStream err:", err)
				c.Status(http.StatusForbidden)
				return
			}
			hash := utils.CalculateSha1(getStream)
			if hash != stream.Hash {
				stream.Commit(false)
				log.Println("resumable put done but hash mismatch")
				log.Println("calculate hash:", hash)
				log.Println("stream.hash:", stream.Hash)
				c.Status(http.StatusForbidden)
				return
			}
			if locate.Exist(url.PathEscape(hash)) {
				stream.Commit(false)
			} else {
				stream.Commit(true)
			}
			file := meta.File{
				Name: stream.Name,
				Hash: stream.Hash,
				Size: strconv.Itoa(int(stream.Size)),
			}
			err = filemeta.Put(file)
			if err != nil {
				log.Println("put fileMeta err:", err)
				c.JSON(http.StatusInternalServerError, response.Err("put fileMeta err"))
				return
			}
			return
		}
	}
}

func GetCurrentSize(c *gin.Context) {
	// todo 获取token方式待定
	token := c.Query("token")
	stream, err := rs.NewRSResumablePutStreamFromToken(token)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusForbidden)
		return
	}
	current := stream.CurrentSize()
	if current == -1 {
		c.Status(http.StatusNotFound)
		return
	}
	c.Header("content-length", fmt.Sprintf("%d", current))
}
