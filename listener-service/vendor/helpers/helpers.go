package helpers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Helpers struct{}

type PayLoad struct{
	Name string `json:"name"`
	Data string `json:"data"`
}

func (h Helpers) readJson(w http.ResponseWriter, r* http.Request, data any) error{
	maxBytes := 1048576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err!= nil{
		return err
	}

	err = dec.Decode(&struct{}{})
	if err!=io.EOF{
		return errors.New("Body Must Have a Single Value")
	}

 	return nil
}

func (h Helpers) writeJson(w http.ResponseWriter, status int, data any, headers ...http.Header) error{
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

func (h Helpers) logEventAsJson(destUrl string, w http.ResponseWriter, load PayLoad){
	jsonData, _ := json.MarshalIndent(load, "","\t")
	request, err := http.NewRequest("POST", destUrl, bytes.NewBuffer(jsonData))
	if err!=nil{
		fmt.Println(err);
		return
	}

	request.Header.Set("Content-Type","application/json")
	client := http.Client{}

	response, err := client.Do(request)
	if err!=nil{
		fmt.Println(err);
		return
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted{
		return
	}
}