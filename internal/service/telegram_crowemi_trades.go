package telegram

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const DEFAULT_ERROR = "Something went wrong. Check my logs for details."

type CrowemiTrades struct {
	BotBase
}

func (c CrowemiTrades) HandleMessage(update Update) {
	// Handle the message for Crowemi Trades Bot
	if c.Config.Crowemi.Debug {
		fmt.Println("Crowemi Trades Bot: Received message:", update)
	}

	botConfig := c.Config.BotConfig[CROWEMI_TRADES]
	botToken := botConfig.Token
	chatID := fmt.Sprintf("%v", update.Message.Chat.ID)

	if c.ValidateUser(int(update.Message.From.ID)) {
		switch update.Message.Text {
		case "/status":

			err := sendMessage(botToken, chatID, fmt.Sprintf("Sure, @%s! Let me check the status for you.", update.Message.From.Username))
			if err != nil {
				fmt.Println("Error sending message:", err)
				return
			}

			client := http.Client{}
			url := fmt.Sprintf("%sstatus/", c.Config.Crowemi.Uri["crowemi-trades"])
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				// TODO: add logging
				fmt.Println("Error creating request:", err)
				_ = sendMessage(botToken, chatID, DEFAULT_ERROR)
				return
			}

			err = c.Config.Crowemi.CreateHeaders(req, c.Config.Crowemi.Uri["crowemi-trades"], "")
			if err != nil {
				// TODO: add logging
				fmt.Println("Error creating headers:", err)
				_ = sendMessage(botToken, chatID, DEFAULT_ERROR)
				return
			}

			resp, err := client.Do(req)
			defer resp.Body.Close()
			if err != nil || resp.StatusCode != http.StatusOK {
				// TODO: add logging
				fmt.Println("Error making request:", err)
				fmt.Println("Response status code:", resp.StatusCode)
				_ = sendMessage(botToken, chatID, DEFAULT_ERROR)
				return
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Error reading response body:", err)
				_ = sendMessage(botToken, chatID, DEFAULT_ERROR)
				return
			}

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
				message := fmt.Sprintf("%s %s: target %.2f; current %.2f; delta %.2f", symbol, key, value.BuyPrice, value.CurrentPrice, value.Diff)
				err := sendMessage(botToken, chatID, message)
				if err != nil {
					fmt.Println("Error sending message:", err)
					_ = sendMessage(botToken, chatID, DEFAULT_ERROR)
					return
				}
			}
			return
		case "/summary":
			return
		default:
			return
		}
	} else {
		_ = sendMessage(botToken, chatID, fmt.Sprintf("%s you are not allowed to use this bot.", update.Message.From.Username))
		return
	}
}
