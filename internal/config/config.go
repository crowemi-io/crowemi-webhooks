package config

import (
	"context"
	"fmt"
	"net/http"

	"google.golang.org/api/idtoken"
)

type crowemi struct {
	ClientName      string            `json:"client_name"`
	ClientId        string            `json:"client_id"`
	ClientSecretKey string            `json:"client_secret_key"`
	Uri             map[string]string `json:"uri"`
	Env             string            `json:"env"`
	Debug           bool              `json:"debug"`
}

func (c crowemi) CreateHeaders(req *http.Request, audience string, sessionID string) error {
	// mirror functionality of crowemi-py-utils
	// TODO: add these to common library
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("crowemi-client-id", c.ClientId)
	req.Header.Set("crowemi-client-secret-key", c.ClientSecretKey)
	req.Header.Set("crowemi-client-name", c.ClientName)
	req.Header.Set("crowemi-session-id", sessionID)
	if c.Env == "dev" || c.Env == "prod" {
		token, err := c.GetAuth(audience)
		if err != nil {
			fmt.Printf("Error getting auth token: %v\n", err)
			return err
		}
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return nil
}
func (c crowemi) GetAuth(targetAudience string) (string, error) {
	ctx := context.Background()

	// Fetch an identity token for the target audience (Cloud Run URL)
	tokenSource, err := idtoken.NewTokenSource(ctx, targetAudience)
	if err != nil {
		fmt.Printf("Error creating token source: %v\n", err)
		return "", err
	}

	token, err := tokenSource.Token()
	if err != nil {
		fmt.Printf("Error retrieving identity token: %v\n", err)
		return "", err
	}

	// Return the identity token
	return token.AccessToken, nil
}

type googleCloudCredential struct{}

type googleCloud struct {
	ProjectId  string                `json:"project_id"`
	Topic      string                `json:"topic"`
	Credential googleCloudCredential `json:"credentials"`
}

type bot struct {
	ChannelId    string `json:"channel_id"`
	Token        string `json:"token"`
	AllowedUsers []int  `json:"allowed_users"`
}

type Webhooks struct {
	App         string         `json:"app"`
	Crowemi     crowemi        `json:"crowemi"`
	BotConfig   map[string]bot `json:"bot"`
	GoogleCloud googleCloud    `json:"google_cloud"`
}
