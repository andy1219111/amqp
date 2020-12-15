# amqp
对streadway/amqp库的二次封装

# 消费者的栗子

```golang
	//建立消费者
	session, err := amqp.NewSession("amqp://work:work@127.0.0.1:5672")
	if err != nil {
		log.Fatalf("connect rabbitMQ failed:%v", err)
		panic("connect rabbitMQ failed")
	}
	defer session.Close()
	//建立消费者
	consumer, err := amqp.NewConsumer(session, "log_queue", "log_upload", "direct", "log_upload")
	if err != nil {
		log.Fatalf("connect channel failed:%v", err)
		panic("connect rabbitMQ failed")
	}
	defer consumer.Close()

	msgs, err := consumer.Consume()
	if err != nil {
		log.Fatalf("connect rabbitMQ failed:%v", err)
		panic("connect rabbitMQ failed")
	}

	for dataBody := range msgs {
		log.Println(string(dataBody.Body))
		dataBody.Ack(false)
	}
  
```

# 生产者的栗子

```golang
	//建立生产者
	session, err := amqp.NewSession("amqp://work:work@192.168.9.153:5672")
	if err != nil {
		log.Fatalf("connect rabbitMQ failed:%v", err)
		panic("connect rabbitMQ failed")
	}
	defer session.Close()

	//建立生产者
	producer, err := amqp.NewProducer(session, "log_upload", "direct", "log_upload")
	if err != nil {
		log.Fatalf("connect channel failed:%v", err)
		panic("connect rabbitMQ failed")
	}
	defer producer.Close()

	err = producer.Publish(
		"log_upload",
		"application/json",
		"this is the body",
	)
	if nil != err {
		log.Fatalf("publish to the mq of the manager server,error:%v", err)
	}
  
```
