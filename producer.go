package amqp

import (
	"github.com/streadway/amqp"
)

//Producer 生产者
type Producer struct {
	ExchangeName string
	ExchangeType string
	Session      *Session
	Ch           *amqp.Channel
}

//NewProducer 得到生产者对象
func NewProducer(session *Session, exchangeName, exchangeType string) (*Producer, error) {
	producer := &Producer{
		Session:      session,
		ExchangeName: exchangeName,
		ExchangeType: exchangeType,
	}
	channel, err := producer.Session.Conn.Channel()
	if err != nil {
		return nil, err
	}
	producer.Ch = channel

	err = producer.Ch.ExchangeDeclare(
		exchangeName, // name of the exchange
		exchangeType, // type
		true,         // durable
		false,        // delete when complete
		false,        // internal
		false,        // noWait
		nil,
	)
	if err != nil {
		return nil, err
	}

	return producer, nil
}

//Publish  消息投递
func (p *Producer) Publish(routingKey, contentType, body string) error {
	err := p.Ch.Publish(
		p.ExchangeName,
		routingKey,
		false, false,
		amqp.Publishing{
			ContentType:  contentType,
			Body:         []byte(body),
			DeliveryMode: amqp.Transient,
		})
	return err
}

//Close  关闭通道
func (p *Producer) Close() error {
	return p.Ch.Close()
}
