package event

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

func ConnectRabbitMQ(rabbitUrl string, retryCount int64)(*amqp.Connection, error){
	var counts int64
	var sleepTime = 2*time.Second
	var connection *amqp.Connection

	for{
		c, err := amqp.Dial(rabbitUrl)
		if err != nil{
			fmt.Printf("Rabbit is Not Ready Sleeping For %d Seconds Counts %d", sleepTime, counts)
			time.Sleep(sleepTime)

			counts++
			sleepTime+=sleepTime*2

		}else {
			fmt.Println("Rabbit is Ready")
			connection = c
			break
		}

		if counts > retryCount{
			fmt.Println("SomeThing went Wrong, quit Trying...")
			break
		}
	}
	return connection, nil
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