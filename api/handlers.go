package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	telegram "github.com/crowemi-io/crowemi-webhooks/internal/service/telegram"
)

func telegramHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ret := fmt.Sprintf(`{"status":"ok", "id": %s }`, r.PathValue("id"))

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error reading body:", err)
		return
	}
	var update telegram.Update
	err = json.Unmarshal(body, &update)
	if err != nil {
		fmt.Println("Error unmarshalling body:", err)
	}
	fmt.Println(update)

	var bot telegram.Bot
	switch r.PathValue("id") {
	case telegram.CROWEMI_TRADES:
		bot = telegram.CrowemiTradesBot{
			Config: Config{},
		}
	default:
		fmt.Println("Unknown bot id")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status":"error", "message":"unknown bot id"}`))
	}

	err = bot.HandleMessage(update)
	if err != nil {
		fmt.Println("Error handling message:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write([]byte(ret))
}
