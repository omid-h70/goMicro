package main

import (
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"logger-service/data"
	"time"
)

type RPCServer struct{}

type RPCPayLoad struct {
	Name string
	Data string
}

var client *mongo.Client

func(r* RPCServer)LogInfo(load RPCPayLoad)(rsp *string){

	data.New(client)
	var logEntry data.LogEntry = data.LogEntry{
		Name:load.Name,
		Data:load.Data,
		CreatedAt: time.Now(),
	}

	err := logEntry.Insert(logEntry)
	logEntry.Insert1()
	if err != nil{
		log.Println(err)
		return
	}
	*rsp = "Got RPC  The Log"
	return nil
}
