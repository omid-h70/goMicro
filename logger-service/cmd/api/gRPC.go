package main

import (
	"context"
	logger_service "helpers/logs"
	"log"
	"logger-service/data"
)

//var client *mongo.Client

type LogServer struct{
	logger_service.UnimplementedLogServiceServer
	models data.Models
}

func (l *LogServer)WriteLog(ctx context.Context, req *logger_service.LogRequest) (*logger_service.LogResponse, error){
	input := req.GetLogEntry()

	logEntry := data.LogEntry{
		Name: input.Name,
		Data: input.Data,
	}

	err := l.models.LogEntry.Insert(logEntry)
	if err != nil{
		log.Println(err)
		return &logger_service.LogResponse{
			Response: "Failed",
		}, err
	}

	logResp := logger_service.LogResponse{
		Response: "gRPC Logging Done",
	}
	return &logResp, nil
}
