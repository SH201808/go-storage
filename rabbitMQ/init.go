package rabbitmq

import (
	"encoding/json"
	"file-server/setting"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

type MQ struct {
	queueName    string
	exchangeName string
	Chan         *amqp.Channel
}

func New(cfg setting.RabbitMQConfig) *MQ {
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		cfg.User,
		cfg.Password,
		cfg.Ip,
		cfg.Port)
	conn, err := amqp.Dial(url)
	failOnError(err, "init mq err")

	c, err := conn.Channel()
	failOnError(err, "get channel err")

	queue, err := c.QueueDeclare(
		"",
		false,
		true,
		false,
		false,
		nil,
	)
	failOnError(err, "declare queue err")

	mq := new(MQ)
	mq.Chan = c
	mq.queueName = queue.Name
	return mq
}

func (mq *MQ) Bind(exchangeName string) {
	err := mq.Chan.QueueBind(
		mq.queueName,
		"",
		exchangeName,
		false,
		nil,
	)
	failOnError(err, "bind exchange error")
	mq.exchangeName = exchangeName
}

func (mq *MQ) Publish(exchangeName string, body interface{}) {
	data, _ := json.Marshal(body)
	err := mq.Chan.Publish(exchangeName,
		"",
		false,
		false,
		amqp.Publishing{
			ReplyTo: mq.queueName,
			Body:    data,
		},
	)
	failOnError(err, "publish msg err")
}

func (mq *MQ) Send(queue string, body interface{}) {
	data, _ := json.Marshal(body)
	err := mq.Chan.Publish("",
		queue,
		false,
		false,
		amqp.Publishing{
			ReplyTo: mq.queueName,
			Body:    data,
		},
	)
	failOnError(err, "send msg err")
}

func (mq *MQ) Consume() <-chan amqp.Delivery {
	c, err := mq.Chan.Consume(mq.queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "consume msg err")
	return c
}

func (mq *MQ) Close() {
	mq.Chan.Close()
}

type QueueInfo struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       map[string]interface{}
}

type PublishInfo struct {
	Exchange  string
	Key       string
	Mandatory bool
	Immediate bool
	Msg       amqp.Publishing
}

type ConsumeInfo struct {
	Queue     string
	Consumer  string
	AutoAck   bool
	Exclusive bool
	NoLocal   bool
	NoWait    bool
	Args      map[string]interface{}
}

func ConstructQueue(Name string, Durable bool, AutoDelete bool, Exclusive bool, NoWait bool, Args map[string]interface{}) *QueueInfo {
	return &QueueInfo{
		Name:       Name,
		Durable:    Durable,
		AutoDelete: AutoDelete,
		Exclusive:  Exclusive,
		NoWait:     NoWait,
		Args:       Args,
	}
}

func ConstructPublish(Exchange string, Key string, Mandatory bool, Immediate bool, data []byte) *PublishInfo {
	return &PublishInfo{
		Exchange:  Exchange,
		Key:       Key,
		Mandatory: Mandatory,
		Immediate: Immediate,
		Msg: amqp.Publishing{
			Body: data,
		},
	}
}

var Conn *amqp.Connection

func Init(config *setting.RabbitMQConfig) *amqp.Connection {
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		config.User,
		config.Password,
		config.Ip,
		config.Port)
	conn, err := amqp.Dial(url)
	failOnError(err, "初始化连接消息队列错误")

	return conn
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
