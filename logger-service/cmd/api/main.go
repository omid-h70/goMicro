package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"helpers"
	loggerService "helpers/logs"
	"log"
	"logger-service/data"
	"net"
	"net/http"
	"net/rpc"
	"time"
)

const(
	webPort = "80"
	rpcPort = "5001"
	gRPCPort = "50001"

	mongoConnTimeout = 10
	mongoURL = "mongodb:://mongo:27017"
	mongoUser = "admin"
	mongoPass = "password"
)

type Config struct{
	Models data.Models
}

func main(){
	client, err := ConnectToMongoDB(mongoURL)
	if err != nil{
		log.Panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), mongoConnTimeout*time.Second)
	defer cancel()

	defer func(){
		client.Disconnect(ctx)
	}()

	app := Config{ data.New(client) }

	rpc.Register(new(RPCServer))
	go rpcListen()

	go gRPCListen()

	loggerAddr := fmt.Sprintf(":%s",helpers.GetEnvVar("LOG_SERVICE_WEB_PORT", webPort));
	server := &http.Server{
		Addr: loggerAddr,
		Handler: app.routes(),
	}

	fmt.Println("Logger Service Has Started On "+loggerAddr)
	log.Panic(server.ListenAndServe())
}

func rpcListen()(error){
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", helpers.GetEnvVar("LOG_SERVICE_RPC_PORT", rpcPort)))
	if err!=nil{
		log.Println(err)
		return err
	}

	for{
		rpcConn, err := listener.Accept()
		if err != nil {
			continue
		}
		go rpc.ServeConn(rpcConn)
	}
}

func gRPCListen()(error){
	loggerAddr := fmt.Sprintf("0.0.0.0:%s", helpers.GetEnvVar("LOG_SERVICE_GRPC_PORT", gRPCPort))
	listener, err := net.Listen("tcp", loggerAddr)
	if err!=nil{
		log.Println(err)
		return err
	}

	s := grpc.NewServer()
	loggerService.RegisterLogServiceServer(s, &LogServer{models:data.New(client)})
	fmt.Println("Logger Service Has Started On "+loggerAddr)

	if err = s.Serve(listener); err != nil{

	}

	defer func(){
		listener.Close()
	}()

	return err
}

func ConnectToMongoDB(dbString string)(*mongo.Client, error){
	clientOptions := options.Client().ApplyURI(dbString)
	clientOptions.SetAuth(options.Credential{
		Username: mongoUser,
		Password: mongoPass,
	})

	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil{
		log.Panic(err)
	}
	return c, nil
}