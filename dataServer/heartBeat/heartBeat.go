package heartBeat

import (
	rabbitmq "file-server/rabbitMQ"
	"file-server/setting"
	"time"
)

func Start() {
	Addr := setting.Conf.MachineIP + setting.Conf.Port
	mq := rabbitmq.New(*setting.Conf.RabbitMQConfig)
	defer mq.Close()

	for {
		mq.Publish("apiServers", Addr)
		time.Sleep(5 * time.Second)
	}
}
