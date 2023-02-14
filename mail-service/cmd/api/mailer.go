package main

import (
	"bytes"
	"fmt"
	mail "github.com/xhit/go-simple-mail/v2"
	"html/template"
	"sync"
	"time"
)

type Mail struct{
	Domain string
	Host string
	Port int
	Username string
	Password string
	Encryption string
	From string
	FromName string
	To string
	FromAddress string

	wait *sync.WaitGroup
	MailerChan chan Message
	ErrorChan chan error
	Done chan bool
}

type Message struct{
	From string
	FromName string
	To string
	Subject string

	Template string
	//Attachment []string
	DataMap map[string]any
}

func(app *App)listenForMail(){
	for{
		select {
			case email := <-app.Mailer.MailerChan:
				app.Mailer.sendMail(email)
			case err := <-app.Mailer.ErrorChan:
				fmt.Println("Something wen wrong", err)
			default:
			case <- app.Mailer.Done:
				return
		}
	}
}

func(app *App)listenForShutDown(){

}

func (m* Mail) sendMail(msg Message){
	defer m.wait.Done()

	if msg.Template == ""{
		msg.Template = "mail"
	}

	if msg.FromName == ""{
		msg.FromName = m.FromAddress
	}

	if msg.From == ""{
		msg.From = m.From
	}

	htmlTextMsg, err := m.buildHTMLMessage(msg)
	if err != nil{
		m.ErrorChan <- err
	}

	plainTextMsg, err := m.buildPlainTextMessage(msg)
	if err != nil{
		m.ErrorChan <- err
	}

	mailServer := mail.NewSMTPClient()
	mailServer.Host = m.Host
	mailServer.Username = m.Username
	mailServer.Password = m.Password
	mailServer.Port = m.Port
	//mailServer.Encryption =
	mailServer.ConnectTimeout = 10 * time.Second
	mailServer.SendTimeout = 10 * time.Second
	mailServer.KeepAlive = false

	mailClient, err := mailServer.Connect()
	if err != nil{
		m.ErrorChan <- err
	}

	email := mail.NewMSG()
	email.SetFrom(m.FromAddress).AddTo(msg.To).SetSubject(msg.Subject)

	email.SetBody(mail.TextPlain, plainTextMsg)
	email.AddAlternative(mail.TextHTML, htmlTextMsg)

	//Add Attachments Later

	err = email.Send(mailClient)
	if err != nil{
		m.ErrorChan <- err
	}
}

func(m* Mail) buildHTMLMessage(msg Message)(string , error){
	templateName := fmt.Sprintf("./cmd/web/%s.gohtml", msg.Template)
	templateToRender, err := template.New("email-html").ParseFiles(templateName)
	if err != nil{

	}

	/*Type bytes.Buffer has implemented io.writer interface*/
	var tpl bytes.Buffer
	if err = templateToRender.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", err
	}
	/*
	formattedMessage, err = m.inlineCSS(formattedMessage)
	if err != nil {
		return "", err
	}
	 */


	plainMessage := tpl.String();
	return plainMessage, nil
}

func(m* Mail) buildPlainTextMessage(msg Message)(string , error){
	templateName := fmt.Sprintf("./cmd/web/%s.gohtml", msg.Template)
	templateToRender, err := template.New("email-plain").ParseFiles(templateName)
	if err != nil{

	}

	/*Type bytes.Buffer has implemented io.writer interface*/
	var tpl bytes.Buffer
	if err = templateToRender.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", err
	}

	plainMessage := tpl.String();
	return plainMessage, nil
}