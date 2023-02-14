package main

import (
	"fmt"
	"helpers/event"
	"log"
	"net/http"
	"os"
)

const (
	rabbitUrl = "amqp://guest:guest@rabbitmq"
	retryCount = 5
    webPort = 80
)

type App struct{
	amqp   event.RabbitMq
	Config AppConfig
}

type ServiceConfig struct{
	Addr string
}

type AppConfig struct{
	LogServiceConfig ServiceConfig
}

func main(){

	rabbitMq := event.RabbitMq{}
	err := rabbitMq.ConnectRabbitMQ(rabbitUrl, retryCount)
	defer rabbitMq.Conn.Close()

	if err != nil{
		fmt.Println(err)
		os.Exit(1)
	}

	app := App{}
	app.Config.LogServiceConfig.Addr = fmt.Sprintf(":%d",webPort);
	server := &http.Server{
		Addr: app.Config.LogServiceConfig.Addr,
		Handler: app.routes(),
	}

	fmt.Println("Broker Service is Started On "+app.Config.LogServiceConfig.Addr)
	log.Panic(server.ListenAndServe())
}
