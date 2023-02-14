package main

import (
	"encoding/json"
	"net/http"
)

type jsonResponse struct{
	Error bool 	   `json:"error"`
	Message string `json:"message"`
	Data any 	   `json:"data,omitempty"` // type "any" is new syntax since 1.18, it means interface{}
}

func (app* App)DoAuth(w http.ResponseWriter, r *http.Request){
	payload := jsonResponse{
		Error: false,
		Message: "Hi Broker",
	}
	out, _ := json.MarshalIndent(payload, "", "\t")
	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Application-Type","application/json")
	w.Write(out)
}

//func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
//	var requestPayload struct {
//		Email    string `json:"email"`
//		Password string `json:"password"`
//	}
//
//	err := app.readJSON(w, r, &requestPayload)
//	if err != nil {
//		app.errorJSON(w, err, http.StatusBadRequest)
//		return
//	}
//
//	// validate the user against the database
//	user, err := app.Models.User.GetByEmail(requestPayload.Email)
//	if err != nil {
//		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
//		return
//	}
//
//	valid, err := user.PasswordMatches(requestPayload.Password)
//	if err != nil || !valid {
//		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
//		return
//	}
//
//	// log authentication
//	err = app.logRequest("authentication", fmt.Sprintf("%s logged in", user.Email))
//	if err != nil {
//		app.errorJSON(w, err)
//		return
//	}
//
//	payload := jsonResponse{
//		Error:   false,
//		Message: fmt.Sprintf("Logged in user %s", user.Email),
//		Data:    user,
//	}
//
//	app.writeJSON(w, http.StatusAccepted, payload)
//}