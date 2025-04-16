package telegram

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/crowemi-io/crowemi-webhooks/internal/config"
)

type CrowemiTrades struct {
	Config config.Webhooks
}

func (c CrowemiTrades) HandleMessage(update Update) error {
	// Handle the message for Crowemi Trades Bot
	switch update.Message.Text {
	case "/status":
		client := http.Client{}
		url := fmt.Sprintf("%sstatus/", c.Config.Crowemi.Uri["crowemi-trades"])
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return fmt.Errorf("error creating request: %v", err)
		}

		// TODO: add these to common library
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("crowemi-client-id", c.Config.Crowemi.ClientId)
		req.Header.Set("crowemi-client-secret-key", c.Config.Crowemi.ClientSecretKey)
		req.Header.Set("crowemi-client-name", c.Config.Crowemi.ClientName)
		req.Header.Set("crowemi-session-id", "crowemi-client-session-id")

		resp, err := client.Do(req)
		defer resp.Body.Close()
		if err != nil || resp.StatusCode != http.StatusOK {
			fmt.Println("Error making request:", err)
			fmt.Println("Response status code:", resp.StatusCode)
			return fmt.Errorf("error making request: %v; status code: %v", err, resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("error reading response body: %v", err)
		}
		fmt.Println("Response body:", string(body))

		status := &StockMap{}
		json.Unmarshal(body, status)
		fmt.Println("Status:", status)

		for key, value := range *status {
			// 🔴 KO: target 73.51; current 73.18; delta -0.33
			var symbol string
			if value.Diff > 0 {
				symbol = "🟢"
			} else {
				symbol = "🔴"
			}
			message := fmt.Sprintf("%s %s: target %f; current %f; delta %f\n", symbol, key, value.BuyPrice, value.CurrentPrice, value.Diff)
			err := sendMessage(c.Config.BotConfig[CROWEMI_TRADES].Token, fmt.Sprintf("%v", update.Message.Chat.ID), message)
			if err != nil {
				fmt.Println("Error sending message:", err)
				return fmt.Errorf("error sending telegram message: %v", err)
			}
		}
		return nil

	case "/summary":
		return nil
	default:
		return nil
	}
}
