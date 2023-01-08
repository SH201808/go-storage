package fileblock

import (
	"context"
	"file-server/dao"
	"file-server/db/redis"
	"file-server/models"
	response "file-server/models/Response"
	"file-server/models/meta"
	"file-server/utils"
	"log"
	"math"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

func InitMeta(c *gin.Context) {
	filehash := c.PostForm("filehash")
	fileSize, _ := strconv.Atoi(c.PostForm("fileSize"))

	upInfo := models.MultipartUploadInfo{
		FileHash:   filehash,
		FileSize:   fileSize,
		UploadId:   strconv.Itoa(int(utils.GenSFID())),
		ChunkSize:  5 * 1024 * 1024,
		ChunkCount: int(math.Ceil(float64(fileSize) / (5 * 1024 * 1024))),
	}
	ctx := context.Background()
	redis.DB.HSet(ctx, "MP_"+upInfo.UploadId, "chunkCount", upInfo.ChunkCount,
		"fileHash", upInfo.FileHash, "fileSize", upInfo.FileSize)

	c.JSON(http.StatusOK, response.Success("初始化信息成功", gin.H{
		"data": upInfo,
	}))
}

func Uppart(c *gin.Context) {
	uploadId := c.PostForm("uploadId")
	chunkIndex := c.PostForm("chunkIndex")
	upPartFile, _ := c.FormFile("upPartFile")

	fpath := "../../data/" + uploadId + "/" + chunkIndex
	err := os.MkdirAll(path.Dir(fpath), 0744)
	if err != nil {
		log.Println("创建文件夹错误：", err)
		c.JSON(http.StatusInternalServerError, response.Err("创建文件夹错误"))
		return
	}

	err = c.SaveUploadedFile(upPartFile, fpath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Err("上传区块错误"))
		return
	}
	//更新redis缓存状态
	ctx := context.Background()
	redis.DB.HSet(ctx, "MP_"+uploadId, "chkidx_"+chunkIndex, 1)

	c.JSON(http.StatusOK, response.Success())
}

func CompleteUpload(c *gin.Context) {
	upid := c.PostForm("uploadId")
	fileHash := c.PostForm("fileHash")
	fileSize, _ := strconv.Atoi(c.PostForm("fileSize"))
	fileName := c.PostForm("fileName")
	userId := c.GetInt("useId")

	redisCtx := context.Background()
	result := redis.DB.HGetAll(redisCtx, "MP_"+upid).Val()

	totalCount := 0
	chunkCount := 0
	for k, v := range result {
		if k == "chunkCount" {
			totalCount, _ = strconv.Atoi(v)
		} else if strings.HasPrefix(k, "chkidx_") && v == "1" {
			chunkCount += 1
		}
	}

	if totalCount != chunkCount {
		c.JSON(http.StatusOK, response.Success("分块上传未完成"))
		return
	}
	completeFile, err := os.Create("../fileStore/" + fileName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Err("创建最终文件错误"))
		return
	}

	//合并
	err = egGoIntegrate(totalCount, completeFile, upid)
	if err != nil {
		log.Println("合并文件错误")
		c.JSON(http.StatusInternalServerError, response.Err("合并文件错误"))
		return
	}

	file := meta.File{
		Sha:      fileHash,
		FileName: fileName,
		FileSize: int64(fileSize),
		Dst:      "../fileStore/" + fileName,
	}
	dao.UploadFileMetaAndUserFile(file, userId)
	c.JSON(http.StatusOK, response.Success())
}

func egGoIntegrate(totalCount int, completeFile *os.File, upid string) error {
	eg, ctx := errgroup.WithContext(context.Background())

	for i := 0; i < totalCount; i++ {
		seq := strconv.Itoa(i)
		offset := int64(i * 5 * 1024 * 1024)
		eg.Go(func() error {
			select {
			case <-ctx.Done():
				return nil
			default:
				file, err := os.Open("../../data/" + upid + "/" + seq)
				if err != nil {
					log.Println("打开分块文件错误")
					return err
				}
				defer file.Close()
				fileInfo, err := file.Stat()
				if err != nil {
					log.Println("获取分块文件信息错误")
					return err
				}

				data := make([]byte, fileInfo.Size())
				_, err = file.Read(data)
				if err != nil {
					log.Println("读取分块文件错误")
					return err
				}

				_, err = completeFile.WriteAt(data, offset)
				if err != nil {
					log.Println("合并文件错误")
					return err
				}
				return nil
			}
		})
	}

	if err := eg.Wait(); err != nil {
		log.Println("合并文件错误", err)
		return err
	}
	return nil
}
