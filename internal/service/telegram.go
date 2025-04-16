package telegram

import (
	"fmt"
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
