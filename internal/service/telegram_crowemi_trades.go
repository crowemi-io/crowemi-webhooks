package telegram

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const DEFAULT_ERROR = "Something went wrong. Check my logs for details."

type CrowemiTrades struct {
	BotBase
}

func (c CrowemiTrades) HandleMessage(update Update) {

	// Handle the message for Crowemi Trades Bot
	botConfig := c.Config.BotConfig[CROWEMI_TRADES]
	botToken := botConfig.Token

	// this determines what chat to respond to
	var chatID string
	// this determines who sent the message
	var fromID int64
	var messageText string

	// what kinds of message are we dealing with?
	if update.Message != nil {
		chatID = fmt.Sprintf("%v", update.Message.Chat.ID)
		fromID = update.Message.From.ID
		messageText = update.Message.Text
	} else if update.ChannelPost != nil {
		chatID = fmt.Sprintf("%v", update.ChannelPost.Chat.ID)
		fromID = update.ChannelPost.SenderChat.ID
		messageText = update.ChannelPost.Text
	} else {
		fmt.Println("Unknown message type")
		return
	}

	input := strings.Split(messageText, " ")
	var command string
	var args []string
	if len(input) > 0 {
		command = input[0]
		if len(input) > 1 {
			args = input[1:]
		}
	} else {
		command = ""
	}

	if c.ValidateMessage(int(fromID)) {
		switch command {
		case "/status":

			err := sendMessage(botToken, chatID, "Sure! Let me check the status for you.")
			if err != nil {
				fmt.Println("Error sending message:", err)
				return
			}

			client := http.Client{}
			var url string
			if len(args) == 1 {
				url = fmt.Sprintf("%sstatus/%s", c.Config.Crowemi.Uri["crowemi-trades"], args[0])
			} else {
				url = fmt.Sprintf("%sstatus/", c.Config.Crowemi.Uri["crowemi-trades"])
			}
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

			for _, value := range *status {
				// 🔴 KO: target 73.51; current 73.18; delta -0.33
				var symbol string
				if value.Diff > 0 {
					symbol = "🟢"
				} else {
					symbol = "🔴"
				}
				message := fmt.Sprintf("%s %s: target %.2f; current %.2f; delta %.2f", symbol, value.Symbol, value.BuyPrice, value.CurrentPrice, value.Diff)
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
			m := "Unknown command. Please try again."
			err := sendMessage(botToken, chatID, m)
			if err != nil {
				fmt.Println(m)
				fmt.Println("Error sending message:", err)
			}
			return
		}
	} else {
		_ = sendMessage(botToken, chatID, "You are not allowed to use this bot.")
		return
	}
}
