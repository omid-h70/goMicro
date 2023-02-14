package event

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

type RabbitMq struct {
	Conn *amqp.Connection
}

func (rmq* RabbitMq)ConnectRabbitMQ(rabbitUrl string, retryCount int64) error{
	var counts int64
	var sleepTime = 2*time.Second

	for{
		c, err := amqp.Dial(rabbitUrl)
		if err != nil{
			fmt.Printf("Rabbit is Not Ready Sleeping For %d Seconds Counts %d", sleepTime, counts)
			time.Sleep(sleepTime)

			counts++
			sleepTime+=sleepTime*2

		}else {
			fmt.Println("Rabbit is Ready")
			rmq.Conn = c
			break
		}

		if counts > retryCount{
			fmt.Println("SomeThing went Wrong, quit Trying...")
			break
		}
	}
	return nil
}

func (rmq* RabbitMq)Push(exchangeName string, event, severity string)error{
	channel, err := rmq.Conn.Channel()
	if err != nil{
		return err
	}

	return channel.PublishWithContext(
		context.Background(),
		exchangeName,
		severity,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body: []byte(event),
		},
		)
}

func (rmq* RabbitMq)Setup() error{
	channel, err := rmq.Conn.Channel()
	if err!= nil{
		return err
	}
	return declareExchange(channel)
}

func declareExchange(channel *amqp.Channel) error{
	return channel.ExchangeDeclare(
		"logs_topic", //name
		"topic", //type
		true, //durable
		false,// autoDelete
		false,//false
		false,//false
		nil,//arguments
	)
}

func declareRandomQueue(channel *amqp.Channel) (amqp.Queue, error){
	return channel.QueueDeclare(
		"", //name
		false, //durable
		false,// autoDelete
		true,//false
		false,//false
		nil,//arguments
	)
}