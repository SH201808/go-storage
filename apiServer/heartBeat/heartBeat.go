package heartbeat

import (
	rabbitmq "file-server/rabbitMQ"
	"file-server/setting"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

var dataServers = make(map[string]time.Time)
var mutex sync.Mutex

func Listen() {
	mq := rabbitmq.New(*setting.Conf.RabbitMQConfig)
	defer mq.Close()
	mq.Bind("apiServers")

	go removeDataServerList()

	c := mq.Consume()
	for msg := range c {
		dataServer, err := strconv.Unquote(string(msg.Body))
		if err != nil {
			panic(err)
		}
		mutex.Lock()
		dataServers[dataServer] = time.Now()
		mutex.Unlock()
	}
}

func removeDataServerList() {
	for {
		time.Sleep(5 * time.Second)
		mutex.Lock()
		for dataServer, t := range dataServers {
			if t.Add(10 * time.Second).Before(time.Now()) {
				delete(dataServers, dataServer)
			}
		}
		mutex.Unlock()
	}
}

func getDataServers() []string {
	mutex.Lock()
	defer mutex.Unlock()

	res := make([]string, len(dataServers))
	i := 0
	for dataServer := range dataServers {
		res[i] = dataServer
		i++
	}
	return res
}

func ChooseDataServers(n int, exclude map[int]string) (ds []string) {
	candidates := make([]string, 0)
	reverseExclude := make(map[string]struct{})
	for _, addr := range exclude {
		reverseExclude[addr] = struct{}{}
	}

	servers := getDataServers()

	for _, server := range servers {
		if _, ok := reverseExclude[server]; !ok {
			candidates = append(candidates, server)
		}
	}

	length := len(candidates)
	if length < n {
		return
	}
	p := rand.Perm(length)
	for i := 0; i < n; i++ {
		ds = append(ds, candidates[p[i]])
	}
	//选择dataServer时正好shutdown
	return
}
