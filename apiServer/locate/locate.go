package locate

import (
	"encoding/json"
	"file-server/models"
	rabbitmq "file-server/rabbitMQ"
	"file-server/rs"
	"file-server/setting"
	"log"
	"time"
)

func FileLoc(fileHash string) (locateInfo map[int]string) {
	mq := rabbitmq.New(*setting.Conf.RabbitMQConfig)

	data := fileHash
	mq.Publish("dataServers", data)

	c := mq.Consume()
	go func() {
		time.Sleep(1 * time.Second)
		mq.Close()
	}()
	locateInfo = make(map[int]string)
	for i := 0; i < rs.ALL_SHARDS; i++ {
		msg := <-c
		if len(msg.Body) == 0 {
			return
		}
		log.Println(string(msg.Body))
		var info models.LocateMessage
		json.Unmarshal(msg.Body, &info)
		locateInfo[info.Id] = info.Addr
	}
	return
}

func Exist(fileHash string) bool {
	return len(FileLoc(fileHash)) >= rs.DATA_SHARDS
}
