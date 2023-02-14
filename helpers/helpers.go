package helpers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
)

type Helpers struct{}

type JSONPayLoad struct{
	Action  string `json:"action"`
	Data  any `json:"data"`
}

type JSONErrorPayLoad struct{
	Status  string `json:"status"`
	Message string `json:"message"`
	Data string `json:"data"`
}

type JSONMailPayLoad struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

type JSONAuthPayLoad struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type JSONLogPayLoad struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func GetEnvVar(key string, defaultValue string) string{
	var val string
	var result bool
	if val, result = os.LookupEnv(key); !result{
		val = defaultValue
	}
	return val
}

func (h Helpers) ReadJson(w http.ResponseWriter, r* http.Request, data any) error{
	maxBytes := 1048576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil{
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF{
		return errors.New("Body Must Have a Single Value")
	}

 	return nil
}

func (h Helpers) WriteJson(w http.ResponseWriter, status int, data any, headers ...http.Header) error{
	out, err := json.Marshal(data);
	if err != nil{
		return err
	}

	if len(headers) > 0 {
		for key,value := range headers[0]{
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, err = w.Write(out)
	if err != nil {
		return err
	}
	return nil
}

func (h Helpers) ErrorJSON(w http.ResponseWriter, data any, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	out, err := json.Marshal(data);
	if err != nil{
	return err
	}

	return h.WriteJson(w, statusCode, out);
}
