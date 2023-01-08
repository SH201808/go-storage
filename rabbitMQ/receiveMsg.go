package rabbitmq

import "github.com/streadway/amqp"

func ReceiveMsg(Conn *amqp.Connection, q *QueueInfo, c *ConsumeInfo, callback func(data []byte)) {
	ch, err := Conn.Channel()
	failOnError(err, "打开通道错误")
	_, err = ch.QueueDeclare(
		q.Name,
		q.Durable,
		q.AutoDelete,
		q.Exclusive,
		q.NoWait,
		q.Args,
	)
	failOnError(err, "创建队列失败")

	msgs, err := ch.Consume(
		c.Queue,
		c.Consumer,
		c.AutoAck,
		c.Exclusive,
		c.NoLocal,
		c.NoWait,
		c.Args,
	)
	failOnError(err, "获取消息失败")

	done := make(chan bool)
	go func() {
		for msg := range msgs {
			callback(msg.Body)
		}
	}()

	<-done
}
