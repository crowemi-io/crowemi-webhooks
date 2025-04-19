package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/crowemi-io/crowemi-webhooks/internal/config"
	telegram "github.com/crowemi-io/crowemi-webhooks/internal/service"
)

type Handlers struct {
	Config *config.Webhooks
}

func (h *Handlers) TelegramHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// its always all good
	w.WriteHeader(http.StatusOK)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error reading body:", err)
		return
	}
	if h.Config.Crowemi.Debug {
		fmt.Println("Request Body:", string(body))
	}

	var update telegram.Update
	err = json.Unmarshal(body, &update)
	if err != nil {
		fmt.Println("Error unmarshalling body:", err)
	}

	var bot telegram.Bot
	switch r.PathValue("id") {
	case telegram.CROWEMI_TRADES:
		fmt.Println("Running crowemi-trades bot")
		bot = telegram.CrowemiTrades{
			BotBase: telegram.BotBase{Config: *h.Config},
		}
	default:
		fmt.Println("Unknown bot id")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status":"error", "message":"unknown bot id"}`))
	}

	bot.HandleMessage(update)
}
