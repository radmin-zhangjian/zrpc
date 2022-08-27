package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type addParam struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type addResult struct {
	Code int `json:"code"`
	Data int `json:"data"`
}

func main() {
	url := "http://127.0.0.1:8090/add"
	param := addParam{
		X: 10,
		Y: 20,
	}
	paramBytes, _ := json.Marshal(param)
	resp, _ := http.Post(url, "application/json", bytes.NewReader(paramBytes))
	defer resp.Body.Close()
	respBytes, _ := ioutil.ReadAll(resp.Body)
	var respData addResult
	json.Unmarshal(respBytes, &respData)
	log.Println(respData)
}
