package telegram

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	CROWEMI_TRADES = "crowemi-trades"
)

type Update struct {
	UpdateID int      `json:"update_id"`
	Message  *Message `json:"message,omitempty"`
}

type Message struct {
	MessageID int    `json:"message_id"`
	From      *User  `json:"from,omitempty"`
	Chat      *Chat  `json:"chat"`
	Date      int64  `json:"date"`
	Text      string `json:"text,omitempty"`
}

type User struct {
	ID           int64  `json:"id"`
	IsBot        bool   `json:"is_bot"`
	FirstName    string `json:"first_name,omitempty"`
	LastName     string `json:"last_name,omitempty"`
	Username     string `json:"username,omitempty"`
	LanguageCode string `json:"language_code,omitempty"`
}

type Chat struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Username  string `json:"username,omitempty"`
	Type      string `json:"type"`
}

type Bot interface {
	HandleMessage(update Update) error
}

type CrowemiTradesBot struct {
	Config Config
}

// TODO: make these inline anonymous structs
type StockData struct {
	BuyPrice     float64 `json:"buy_price"`
	CurrentPrice float64 `json:"current_price"`
	Diff         float64 `json:"diff"`
}
type StockMap map[string]StockData

func sendMessage(bot_id string, channel_id string, message string) error {
	// uri = "https://api.telegram.org/bot{bot_id}/sendMessage?chat_id={channel_id}&text={message}"
	client := http.Client{}
	uri := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%s&text=%s", bot_id, channel_id, url.QueryEscape(message))
	req, err := http.NewRequest("POST", uri, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
	}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Println("Error making request:", err)
	}
	return nil
}

func (c CrowemiTradesBot) HandleMessage(update Update) error {
	// Handle the message for Crowemi Trades Bot
	switch update.Message.Text {
	case "/status":
		client := http.Client{}
		url := fmt.Sprintf("%sstatus/", c.Config.Crowemi.Uri["crowemi-trades"])
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			// TOOD: handle error
			println("Error creating request:", err)
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
			return err
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
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
			sendMessage(c.Config.BotConfig[0].Token, "", message)
		}
		return nil

	case "/summary":
		return nil
	default:
		return nil
	}
	return nil
}
