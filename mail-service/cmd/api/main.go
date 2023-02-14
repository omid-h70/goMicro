package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

type App struct{
	Mailer Mail
}

func main(){
	app := App{
		Mailer: createMail(),
	}

	webPort := os.Getenv("SERVER_MAIL_PORT")
	log.Println("Starting mail service on port", webPort)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.makeRoutes(),
	}

	go app.listenForMail()
	go app.listenShutDown()

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func createMail() Mail {

	port, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))
	m := Mail{
		Domain:      os.Getenv("MAIL_DOMAIN"),
		Host:        os.Getenv("MAIL_HOST"),
		Port:        port,
		Username:    os.Getenv("MAIL_USERNAME"),
		Password:    os.Getenv("MAIL_PASSWORD"),
		Encryption:  os.Getenv("MAIL_ENCRYPTION"),
		FromName:    os.Getenv("FROM_NAME"),
		FromAddress: os.Getenv("FROM_ADDRESS"),

		Done:		 make(chan bool),
		ErrorChan:	 make(chan error),
		MailerChan:	 make(chan Message, 100), //Buffered Channel
	}

	return m
}