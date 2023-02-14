package main

import (
	"fmt"
	"helpers/event"
	"os"
)

//AMQP stands for Advanced Queueing Message Protocol

const rabbitUrl string = "amqp://guest:guest@rabbitmq"
const retryCount = 5

func main(){

	rabbitMqConn, err := event.ConnectRabbitMQ(rabbitUrl, retryCount)
	defer rabbitMqConn.Close()
	if err != nil{
		fmt.Println(err)
		os.Exit(1)
	}

	consumer, err := event.NewConsumer(rabbitMqConn)
	if err != nil{
		panic(err)
	}

	err = consumer.Listen([]string{"log.ERROR", "log.WARNING", "log.ERROR"})
	if err != nil {
		fmt.Println(err)
	}
}
