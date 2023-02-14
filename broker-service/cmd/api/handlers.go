package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"helpers"
	logger_service "helpers/logs"
	"net/http"
	"net/rpc"
	"time"
)

func (app* App)Broker(w http.ResponseWriter, r *http.Request){
	payload := helpers.JSONPayLoad{
		Action: "log",
		Data: "Hi Broker",
	}
	out, _ := json.MarshalIndent(payload, "", "\t")
	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Application-Type","application/json")
	w.Write(out)
}

func (app* App)logEventsByRabbit(w http.ResponseWriter, load helpers.JSONPayLoad) error{
	return app.pushToQueue(load)
}

func (app *App)pushToQueue(log helpers.JSONPayLoad) error{
	jData,_ := json.Marshal(log.Data)
	return app.amqp.Push("logs.INFO", string(jData), "dude")
}

func (app *App)logEventAsJson(destUrl string, w http.ResponseWriter, data any){
	jsonData, _ := json.MarshalIndent(data, "","\t")
	request, err := http.NewRequest("POST", destUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println(err);
		return
	}

	request.Header.Set("Content-Type","application/json")
	client := http.Client{}

	response, err := client.Do(request)
	if err != nil{
		fmt.Println(err);
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted{
		return
	}
}

func (app *App)logEventByRPC(destUrl string, w http.ResponseWriter, data any){
	var h helpers.Helpers
	client, err := rpc.Dial("tcp", "logger-service:5001")
	if err != nil {
		fmt.Println(err);
		return
	}

	var result string
	err = client.Call("RPCServer.LogInfo", data, &result)
	if err != nil {
		fmt.Println(err);
		return
	}

	h.WriteJson(w, http.StatusAccepted, data)
}

func (app *App)logEventByGRPC(destUrl string, w http.ResponseWriter, data any){
	var h helpers.Helpers
	conn, err := grpc.Dial("logger-service:50001", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		fmt.Println(err);
		return
	}
	defer conn.Close()

	client := logger_service.NewLogServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err = client.WriteLog(ctx, &logger_service.LogRequest{
		LogEntry: &logger_service.Log{
			Name: "meh",
			Data: "test",
		},
	})
	if err != nil {
		fmt.Println(err);
		return
	}

	h.WriteJson(w, http.StatusAccepted, data)
}

// HandleSubmission is the main point of entry into the broker. It accepts a JSON
// payload and performs an action based on the value of "action" in that JSON.
// it Handles Requests From FrontEnd as well
func (app *App) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload helpers.JSONPayLoad
	var h helpers.Helpers

	err := h.ReadJson(w, r, &requestPayload)
	if err != nil {
		h.ErrorJSON(w, err.Error())
		return
	}

	switch requestPayload.Action {
	case "auth":
		//app.authenticate(w, requestPayload.Auth)
	case "json-log":
		app.logEventAsJson(app.Config.LogServiceConfig.Addr, w, requestPayload)
	case "log":
		app.logEventsByRabbit(w, requestPayload)
	case "mail":
		//app.sendMail(w, requestPayload.Mail)
	default:
		//app.errorJSON(w, errors.New("unknown action"))
	}
}