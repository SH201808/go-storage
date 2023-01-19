package file

import (
	heartbeat "file-server/apiServer/heartBeat"
	"file-server/apiServer/locate"
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

func Upload(c *gin.Context) {
	hash := c.Request.Header.Get("Digests")
	if locate.Exist(hash) {
		c.JSON(http.StatusOK, response.Success("save file success"))
		return
	}
	// c.Request.URL.EscapedPath()
	// todo 上传携带空格的数据存在转义问题
	FileSize := c.Request.Header.Get("FileSize")
	fileSize, _ := strconv.Atoi(FileSize)
	//post元数据到数据服务
	err := storeObject(c.Request.Body, int64(fileSize), hash)
	if err != nil {
		log.Println("Upload err:", err)
		c.JSON(http.StatusInternalServerError, response.Err("upload file err"))
		return
	}

	name := c.Query("fileName")
	file := meta.File{
		Name: name,
		Hash: hash,
		Size: c.Request.Header.Get("FileSize"),
	}
	err = filemeta.Put(file)
	if err != nil {
		log.Println("put fileMeta err:", err)
		c.JSON(http.StatusInternalServerError, response.Err("put fileMeta err"))
		return
	}
	c.JSON(http.StatusOK, response.Success("upload file success"))
}

func storeObject(r io.Reader, fileSize int64, fileHash string) error {
	stream, err := putStream(fileHash, fileSize)
	if err != nil {
		return err
	}

	reader := io.TeeReader(r, stream)
	getHash := utils.CanculateSha1(reader)
	if getHash != fileHash {
		stream.Commit(false)
		return fmt.Errorf("object hash mismatch")
	}
	stream.Commit(true)
	return nil
}

func putStream(fileHash string, fileSize int64) (*rs.RSPutStream, error) {
	servers := heartbeat.ChooseDataServers(rs.ALL_SHARDS, nil)
	if len(servers) != rs.ALL_SHARDS {
		return nil, fmt.Errorf("cannot find enough dataServer")
	}

	return rs.NewRsPutStream(servers, fileHash, fileSize)
}

func Download(c *gin.Context) {
	fileName := c.Query("fileName")
	fileVersion, _ := strconv.Atoi(c.Query("fileVersion"))

	file, err := filemeta.Get(fileName, fileVersion)
	if err != nil {
		log.Println("get fileMeta err: ", err)
		c.JSON(http.StatusInternalServerError, response.Err("get fileMeta error"))
		return
	}
	if file.Hash == "" {
		log.Printf("get hash nil, fileName: %s, fileVersion:%d\n", fileName, fileVersion)
		c.JSON(http.StatusOK, response.Success("hash nil"))
		return
	}
	size, _ := strconv.Atoi(file.Size)
	stream, err := GetStream(file.Hash, int64(size))
	if err != nil {
		log.Println("getStream err:", err)
		c.JSON(http.StatusInternalServerError, response.Err("getStream err"))
		return
	}
	defer stream.Close()

	//获取断点下载偏移量
	offset := utils.GetOffsetFromHeader(c.Request.Header)
	if offset != 0 {
		stream.Seek(offset, io.SeekCurrent)
		c.Header("content-range", fmt.Sprintf("bytes%d-%d/%d", offset, size-1, size))
		c.Status(http.StatusPartialContent)
	}

	_, err = io.Copy(c.Writer, stream)
	if err != nil {
		log.Println("copy err:", err)
		c.JSON(http.StatusInternalServerError, response.Err("copy err"))
		return
	}

	// ip := locate.FileLoc(file.Hash)

	// proxy := httputil.ReverseProxy{
	// 	Director: func(req *http.Request) {
	// 		deleteQuery(c.Request, "fileName", "fileVersion")
	// 		AddQuery(c.Request, "fileHash", file.Hash)
	// 		GenProxy(ip, c.Request.Method, c, req)
	// 	},
	// }

	// proxy.ServeHTTP(c.Writer, c.Request)
}

func GetStream(hash string, size int64) (*rs.RSGetStream, error) {
	locateinfo := locate.FileLoc(hash)
	if len(locateinfo) < rs.DATA_SHARDS {
		return nil, fmt.Errorf("object %s locate fail, result %v", hash, locateinfo)
	}
	dataServers := make([]string, 0)
	if len(locateinfo) != rs.ALL_SHARDS {
		dataServers = heartbeat.ChooseDataServers(rs.ALL_SHARDS-len(locateinfo), locateinfo)
	}
	return rs.NewRSGetStream(locateinfo, dataServers, hash, size)
}

func GenProxy(ip string, method string, c *gin.Context, req *http.Request) {
	url, err := url.Parse("http://" + ip + ":8080" + c.Request.URL.Path)
	if err != nil {
		log.Println("parse url err: ", err)
		return
	}
	req.URL = url
	req.Host = url.Host
	req.URL.RawQuery = c.Request.URL.Query().Encode()
	req.URL.Scheme = "http"
	req.Method = c.Request.Method
}
