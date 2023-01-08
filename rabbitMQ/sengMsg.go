package rabbitmq

import "github.com/streadway/amqp"

func SendMsg(Conn *amqp.Connection, q *QueueInfo, p *PublishInfo) error {
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
	err = ch.Publish(
		p.Exchange,
		p.Key,
		p.Mandatory,
		p.Immediate,
		p.Msg,
	)
	failOnError(err, "发送消息失败")
	return nil
}
