package locate

import (
	"file-server/models"
	rabbitmq "file-server/rabbitMQ"
	"file-server/setting"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

var (
	FileLoc, TempLoc string
	mutex            sync.Mutex
	store            = make(map[string]int, 0)
)

func SetFileLoc() {
	var port string
	flag.StringVar(&port, "port", ":8080", "端口号")
	flag.Parse()
	setting.Conf.Port = port
	FileLoc = "../fileStore/" + setting.Conf.Port + "/"
	TempLoc = "../tempData/" + setting.Conf.Port + "/"
	mkdir(FileLoc)
	mkdir(TempLoc)
}

func mkdir(loc string) {
	_, err := os.Stat(loc)
	if os.IsExist(err) {
		return
	}
	os.Mkdir(loc, os.ModePerm)
	os.Chmod(loc, 0777)
}

func Start() {
	mq := rabbitmq.New(*setting.Conf.RabbitMQConfig)
	getAllFiles()
	mq.Bind("dataServers")

	c := mq.Consume()
	log.Println("消费")
	for msg := range c {
		log.Println("msgBody:", string(msg.Body))
		// log.Println("hash:", string(msg.Body)[7:])
		fileHash, _ := strconv.Unquote(string(msg.Body))
		id := IsExistsFile(fileHash)
		log.Println("id:", id)
		if id != -1 {
			replyMsg := models.LocateMessage{
				Addr: setting.Conf.MachineIP + setting.Conf.Port,
				Id:   id,
			}
			mq.Send(msg.ReplyTo, replyMsg)
		}
	}
}

func getAllFiles() {
	dir, err := ioutil.ReadDir(FileLoc)
	if err != nil {
		log.Fatalln("open dir err:", err)
	}
	log.Println("dir length:", len(dir))
	for _, file := range dir {
		fileName := file.Name()
		log.Println("fileName:", fileName)
		temp := strings.Split(fileName, ".")
		hash := temp[0]
		log.Println("hash:", hash)
		id, _ := strconv.Atoi(temp[1])
		log.Println("id:", id)
		store[hash] = id
	}
}

func IsExistsFile(fileName string) int {
	log.Println("exist fileName:", fileName)
	mutex.Lock()
	defer mutex.Unlock()
	id, ok := store[fileName]
	log.Println("exist id:", id)
	if !ok {
		return -1
	}
	return id
}

func Add(fileHash string, id int) {
	mutex.Lock()
	store[fileHash] = id
	mutex.Unlock()
}

func Delete(fileLoc string) {
	mutex.Lock()
	delete(store, fileLoc)
	mutex.Unlock()
}
