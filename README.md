# amqp
对streadway/amqp库的二次封装

#消费者的栗子

```

//建立消费者
	consumer, err := uamqp.NewConsumer(
		configObj.Amqp["url"],
		configObj.Amqp["exchange"],
		configObj.Amqp["exchange_type"],
		configObj.Amqp["queue_name"],
		configObj.Amqp["router_key"],
	)
	defer consumer.Close()

	msgs, err := consumer.Consume()
	if err != nil {
		logger.Fatalf("connect rabbitMQ failed:%v", err)
		panic("connect rabbitMQ failed")
	}
  
  for dataBody := range msgs {
    fmt.Print(dataBody.Body)
    dataBody.Ack(false)
  }
```
