package main

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"

	"github.com/crowemi-io/crowemi-webhooks/api"
	webhooks "github.com/crowemi-io/crowemi-webhooks/internal/config"
)

func config() webhooks.Webhooks {
	value := os.Getenv("CONFIG")
	decode, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		println("Error decoding config:", err)
	}
	c := webhooks.Webhooks{}
	err = json.Unmarshal(decode, &c)
	if err != nil {
		println("Error unmarshalling config:", err)
	}
	return c
}

func main() {
	config := config()
	handlers := api.Handlers{
		Config: &config,
	}

	http.HandleFunc("/v1/telegram/{id}", func(w http.ResponseWriter, r *http.Request) { handlers.TelegramHandler(w, r) })
	http.ListenAndServe(":8003", nil)
}
