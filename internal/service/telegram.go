package telegram

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/crowemi-io/crowemi-webhooks/internal/config"
)

const (
	CROWEMI_TRADES = "crowemi-trades"
)

type Update struct {
	UpdateID    int          `json:"update_id"`
	Message     *Message     `json:"message,omitempty"`
	ChannelPost *ChannelPost `json:"channel_post,omitempty"`
}

type ChannelPost struct {
	MessageID  int    `json:"message_id"`
	SenderChat Chat   `json:"sender_chat"`
	Chat       Chat   `json:"chat"`
	Date       int64  `json:"date"`
	Text       string `json:"text,omitempty"`
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
	Title     string `json:"title,omitempty"`
	Type      string `json:"type"`
}

type Bot interface {
	HandleMessage(update Update)
}

type BotBase struct {
	Config config.Webhooks
}

func (b BotBase) ValidateMessage(ID int) bool {
	// Validate user ID
	for _, id := range b.Config.BotConfig[CROWEMI_TRADES].AllowedUsers {
		if id == ID {
			return true
		}
	}
	for _, id := range b.Config.BotConfig[CROWEMI_TRADES].AllowedChats {
		if id == ID {
			return true
		}
	}
	return false
}

// TODO: make these inline anonymous structs
// TODO: these need to GO!
type StockData struct {
	Symbol       string  `json:"symbol"`
	BuyPrice     float64 `json:"buy_price"`
	CurrentPrice float64 `json:"current_price"`
	Diff         float64 `json:"diff"`
}
type StockMap []StockData

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
