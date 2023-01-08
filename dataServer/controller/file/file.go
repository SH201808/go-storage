package file

import (
	"crypto/sha256"
	"encoding/base64"
	"file-server/dataServer/locate"
	response "file-server/models/Response"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func Upload(c *gin.Context) {
	fileName := c.PostForm("fileName")
	log.Println(fileName)

	file, err := c.FormFile("file")
	if err != nil {
		log.Println("form file err: ", err)
		c.JSON(http.StatusInternalServerError, response.Err("dataServer form file err"))
		return
	}
	dst := "../../dataServer/fileStore/" + c.Request.Header.Get("Digests")
	c.SaveUploadedFile(file, dst)
	c.JSON(http.StatusOK, response.Success("save file success"))
}

func Download(c *gin.Context) {
	fileHash := c.Query("fileHash")
	filePath := getFile(fileHash)
	c.File(filePath)
}

func getFile(hash string) string {
	files, _ := filepath.Glob(locate.FileLoc + hash + ".*")
	if len(files) != 1 {
		return ""
	}
	file := files[0]
	h := sha256.New()
	sendFile(h, file)
	d := url.PathEscape(base64.StdEncoding.EncodeToString(h.Sum(nil)))
	if d != hash {
		log.Println("object hash mismatch", file)
		locate.Delete(file)
		return ""
	}
	return file
}

func sendFile(w io.Writer, filePath string) {
	file, _ := os.Open(filePath)
	defer file.Close()
	io.Copy(w, file)
}
