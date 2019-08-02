package amqp

import (
	"github.com/streadway/amqp"
)

type Consumer struct {
	uri          string
	exchangeName string
	exchangeType string
	queueName    string
	routingKey   string
	Session      *Session
	Queue        amqp.Queue
}

type Session struct {
	ch      *amqp.Channel
	conn    *amqp.Connection
	errChan chan *amqp.Error
}

func NewConsumer(uri, exchange, exchangeType, queue, routingKey string) (*Consumer, error) {
	consumer := &Consumer{
		uri:          uri,
		exchangeName: exchange,
		exchangeType: exchangeType,
		queueName:    queue,
		routingKey:   routingKey,
	}
	session, err := consumer.Connect()
	if err != nil {
		return nil, err
	}

	consumer.Session = session
	err = session.ch.ExchangeDeclare(
		consumer.exchangeName, // name of the exchange
		consumer.exchangeType, // type
		true,                  // durable
		false,                 // delete when complete
		false,                 // internal
		false,                 // noWait
		nil,
	)
	if err != nil {
		return nil, err
	}

	consumer.Queue, err = session.ch.QueueDeclare(
		queue, // name
		true,  // durable  持久性的,如果事前已经声明了该队列，不能重复声明
		false, // delete when unused
		false, // exclusive 如果是真，连接一断开，队列删除
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return consumer, err
	}
	// 队列和交换机绑定，即是队列订阅了发到这个交换机的消息
	err = session.ch.QueueBind(
		consumer.Queue.Name,
		routingKey,
		exchange,
		false,
		nil)
	if err != nil {
		return nil, err
	}

	return consumer, nil
}

func (c *Consumer) Consume() (<-chan amqp.Delivery, error) {
	msg, err := c.Session.ch.Consume(
		c.queueName, // queue
		"",          // consumer
		false,       // auto-ack   设置为真自动确认消息
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args)
	)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func (c *Consumer) Connect() (*Session, error) {
	errChan := make(chan *amqp.Error)
	s := &Session{errChan: errChan}
	conn, err := amqp.Dial(c.uri)
	if err != nil {
		return nil, err
	}
	go func() {
		errs := conn.NotifyClose(make(chan *amqp.Error))
		for err := range errs {
			s.errChan <- err
		}
	}()

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	s.conn = conn
	s.ch = channel
	return s, err
}

func (c *Consumer) Close() {
	c.Session.ch.Close()
	c.Session.conn.Close()
}
