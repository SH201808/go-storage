package main

import (
	"encoding/json"
	ossImplement "file-server/ossImplement"
	rabbitmq "file-server/rabbitMQ"
	"file-server/setting"
	"fmt"
	"log"
	"math"
	"os"
	"time"

	oss "github.com/aliyun/aliyun-oss-go-sdk/oss"
)

const (
	MinSizeForBlock = 10 * 1024 * 1024
	BlockSize       = 10 * 1024 * 1024
)

func main() {
	q := &rabbitmq.QueueInfo{
		Name:       "toOss",
		Durable:    false,
		AutoDelete: false,
		Exclusive:  false,
		NoWait:     false,
		Args:       nil,
	}
	c := &rabbitmq.ConsumeInfo{
		Queue:     q.Name,
		Consumer:  "",
		AutoAck:   false,
		Exclusive: false,
		NoLocal:   false,
		NoWait:    false,
		Args:      nil,
	}
	setting.Init()
	Conn := rabbitmq.Init(setting.Conf.RabbitMQConfig)
	ossImplement.Init(setting.Conf.OssConfig)
	log.Println("等待数据中")
	rabbitmq.ReceiveMsg(Conn, q, c, transfer)
}

func transfer(data []byte) {
	fileMeta := new(ossImplement.FileMeta)
	json.Unmarshal(data, &fileMeta)

	var err error
	if fileMeta.FileSize < MinSizeForBlock {
		err = DirectUpload(setting.Conf.OssConfig, fileMeta.AbsLoc)
	} else {
		chunckCount := int(math.Ceil(float64(fileMeta.FileSize) / (10 * 1024 * 1024)))
		err = BlockUpload(setting.Conf.OssConfig, fileMeta.AbsLoc, chunckCount)
	}

	if err != nil {
		log.Println("发送到oss错误:", err)
		return
	}

}

func DirectUpload(ossConfig *setting.OssConfig, fileLoc string) error {
	// 填写存储空间名称，例如examplebucket。
	bucket, err := ossImplement.Client.Bucket(ossConfig.BucketName)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	// 依次填写Object的完整路径（例如exampledir/exampleobject.txt）和本地文件的完整路径（例如D:\\localpath\\examplefile.txt）。
	err = bucket.PutObjectFromFile("oss://file-server-sh/"+fileLoc, fileLoc)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	return nil
}

func BlockUpload(ossConfig *setting.OssConfig, fileLoc string, chunckCount int) error {
	// 填写存储空间名称。
	bucketName := ossConfig.BucketName
	// 填写Object完整路径。Object完整路径中不能包含Bucket名称。
	objectName := "oss://file-server-sh/" + fileLoc
	// 填写本地文件的完整路径。如果未指定本地路径，则默认从示例程序所属项目对应本地路径中上传文件。
	locaFilename := fileLoc

	// 获取存储空间。
	bucket, err := ossImplement.Client.Bucket(bucketName)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	// 将本地文件分片，且分片数量指定为3。
	chunks, err := oss.SplitFileByPartNum(locaFilename, chunckCount)
	fd, err := os.Open(locaFilename)
	defer fd.Close()

	// 指定过期时间。
	expires := time.Date(2049, time.January, 10, 23, 0, 0, 0, time.UTC)
	// 如果需要在初始化分片时设置请求头，请参考以下示例代码。
	options := []oss.Option{
		oss.MetadataDirective(oss.MetaReplace),
		oss.Expires(expires),
		// 指定该Object被下载时的网页缓存行为。
		// oss.CacheControl("no-cache"),
		// 指定该Object被下载时的名称。
		// oss.ContentDisposition("attachment;filename=FileName.txt"),
		// 指定该Object的内容编码格式。
		// oss.ContentEncoding("gzip"),
		// 指定对返回的Key进行编码，目前支持URL编码。
		// oss.EncodingType("url"),
		// 指定Object的存储类型。
		// oss.ObjectStorageClass(oss.StorageStandard),
	}

	// 步骤1：初始化一个分片上传事件，并指定存储类型为标准存储。
	imur, err := bucket.InitiateMultipartUpload(objectName, options...)
	// 步骤2：上传分片。
	var parts []oss.UploadPart
	for _, chunk := range chunks {
		fd.Seek(chunk.Offset, os.SEEK_SET)
		// 调用UploadPart方法上传每个分片。
		part, err := bucket.UploadPart(imur, fd, chunk.Size, chunk.Number)
		if err != nil {
			fmt.Println("Error:", err)
			return err
		}
		parts = append(parts, part)
	}

	// 指定Object的读写权限为公共读，默认为继承Bucket的读写权限。
	objectAcl := oss.ObjectACL(oss.ACLPublicRead)

	// 步骤3：完成分片上传，指定文件读写权限为公共读。
	cmur, err := bucket.CompleteMultipartUpload(imur, parts, objectAcl)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	fmt.Println("cmur:", cmur)
	return nil
}
