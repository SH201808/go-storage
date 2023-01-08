package tempFile

import (
	"encoding/json"
	"file-server/dataServer/UUID"
	"file-server/dataServer/locate"
	"file-server/models"
	response "file-server/models/Response"
	"file-server/utils"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

func UploadMeta(c *gin.Context) {
	fileMeta := models.TempFileMeta{
		Size: c.Request.Header.Get("Size"),
		Name: c.Request.Header.Get("Hash"),
		UUID: UUID.Gen(),
	}
	data, err := json.Marshal(fileMeta)
	if err != nil {
		log.Println("Marshal file err:", err)
		c.JSON(http.StatusInternalServerError, response.Err("Craete file error"))
		return
	}

	infoPath := locate.TempLoc + fileMeta.UUID
	infoFile, err := os.OpenFile(infoPath, os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		log.Println("Create infoFile err:", err)
		c.JSON(http.StatusInternalServerError, response.Err("Craete file error"))
		return
	}
	defer infoFile.Close()

	_, err = infoFile.Write(data)
	if err != nil {
		log.Println("Write file err:", err)
		c.JSON(http.StatusInternalServerError, response.Err("Craete file error"))
		return
	}
	datPath := infoPath + ".dat"
	if _, err = os.OpenFile(datPath, os.O_CREATE, 0777); err != nil {
		log.Println("Create datFile err:", err)
		c.JSON(http.StatusInternalServerError, response.Err("Craete file error"))
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{"uuid": fileMeta.UUID}))
}

func UploadtoTempFile(c *gin.Context) {
	uuid := c.Request.Header.Get("uuid")
	tempInfo, err := readFromFile(uuid)
	if err != nil {
		log.Println("read from file err:", err)
		return
	}

	infoPath := locate.TempLoc + tempInfo.UUID
	datPath := infoPath + ".dat"
	datFile, err := os.OpenFile(datPath, os.O_RDWR, 0777)
	if err != nil {
		log.Println("open datFile err:", err)
		c.JSON(http.StatusInternalServerError, response.Err("upload file err"))
		return
	}
	defer datFile.Close()
	_, err = io.Copy(datFile, c.Request.Body)
	if err != nil {
		log.Println("get datFile err:", err)
		c.JSON(http.StatusInternalServerError, response.Err("upload file err"))
		return
	}
	if !compareSize(tempInfo, datFile) {
		os.Remove(infoPath)
		os.Remove(datPath)
		c.JSON(http.StatusOK, response.Err("upload file err"))
		return
	}

	c.JSON(http.StatusOK, response.Success())
}

func DeleteFile(c *gin.Context) {
	uuid := c.Request.Header.Get("uuid")
	infoPath := locate.TempLoc + uuid
	err := os.Remove(infoPath)
	if err != nil {
		log.Println("remove file err: ", err)
		return
	}

	datPath := infoPath + ".dat"
	err = os.Remove(datPath)
	if err != nil {
		log.Println("remove file err: ", err)
		return
	}
}

func RemoveToStore(c *gin.Context) {
	uuid := c.Request.Header.Get("uuid")
	tempInfo, err := readFromFile(uuid)
	if err != nil {
		log.Println("read from file err:", err)
		return
	}

	infoPath := locate.TempLoc + tempInfo.UUID
	datPath := infoPath + ".dat"
	datFile, err := os.Open(datPath)
	if err != nil {
		log.Println("open datFile err:", err)
		c.JSON(http.StatusInternalServerError, response.Err("upload file err"))
		return
	}
	defer datFile.Close()

	os.Remove(infoPath)
	if !compareSize(tempInfo, datFile) {
		os.Remove(datPath)
		c.JSON(http.StatusOK, response.Err("upload file err"))
		return
	}
	commitTempObject(datPath, tempInfo)
}

func commitTempObject(datFile string, tempInfo *models.TempFileMeta) {
	f, _ := os.Open(datFile)
	d := url.PathEscape(utils.CanculateSha1(f))
	f.Close()

	hash := tempInfo.Hash()
	newPath := locate.FileLoc + tempInfo.Name + "." + d
	id, _ := strconv.Atoi(tempInfo.Id())
	os.Rename(datFile, newPath)
	locate.Add(hash, id)
}

func readFromFile(uuid string) (*models.TempFileMeta, error) {
	path := locate.TempLoc + uuid
	f, err := os.OpenFile(path, os.O_RDWR, 0777)
	if err != nil {
		return nil, fmt.Errorf("open file err:", err)
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("read file err:", err)
	}
	var fileMeta models.TempFileMeta
	json.Unmarshal(data, &fileMeta)
	return &fileMeta, nil
}

func compareSize(tempInfo *models.TempFileMeta, datFile *os.File) bool {
	info, err := datFile.Stat()
	if err != nil {
		log.Println("get file stat err:", err)
		return false
	}

	datSize := strconv.Itoa(int(info.Size()))
	if datSize != tempInfo.Size {
		log.Println("file Size mismatch: actualSize:" + datSize + "tempInfoSize:" + tempInfo.Size)
		return false
	}
	return true
}
