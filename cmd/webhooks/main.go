package main

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"
)

func config() Config {
	value := os.Getenv("CONFIG")
	decode, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		println("Error decoding config:", err)
	}
	c := Config{}
	err = json.Unmarshal(decode, &c)
	if err != nil {
		println("Error unmarshalling config:", err)
	}
	return c
}

func main() {
	config := config()
	println("Client Name: ", config.Crowemi.ClientName)

	http.HandleFunc("POST /v1/telegram/{id}", telegramHandler)
	http.ListenAndServe(":8003", nil)
}
