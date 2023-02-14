package main

import (
	"encoding/json"
	"helpers"
	"logger-service/data"
	"net/http"
)

func (app* Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	var requestPayload helpers.JSONPayLoad
	var h helpers.Helpers

	err := h.ReadJson(w, r, &requestPayload)
	if err != nil {
		err = h.ErrorJSON(w, err.Error())
		return
	}

	logPayLoad := data.LogEntry{
		Name: requestPayload.Action,
		//Data: requestPayload.Data,
	}

	err = app.Models.LogEntry.Insert(logPayLoad)
	if err != nil {
		err = h.ErrorJSON(w, err.Error())
		return
	}

	jsonResp := helpers.JSONPayLoad{
		Action: "Log",
		Data:"Accepted",
	}
	jData,_ := json.Marshal(jsonResp)
	h.WriteJson(w, http.StatusAccepted, jData)
}
