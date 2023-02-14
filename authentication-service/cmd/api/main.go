package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

//const webPort = 8081
type App struct{}

func main(){
	app := App{}
	brokerAddr := fmt.Sprintf(":%s",os.Getenv("AUTH_SERVICE_HOST_PORT"));
	server := &http.Server{
		Addr: brokerAddr,
		Handler: app.routes(),
	}

	fmt.Println("Authentication Service is Started On "+brokerAddr)
	log.Panic(server.ListenAndServe())
}
