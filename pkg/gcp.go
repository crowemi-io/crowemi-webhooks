package pkg

import (
	"context"
	"fmt"

	"golang.org/x/oauth2/google"
)

func GetAuth() (string, error) {

	ctx := context.Background()
	ret := ""

	// Fetch default credentials
	creds, err := google.FindDefaultCredentials(ctx)
	if err != nil {
		fmt.Printf("Error finding default credentials: %v\n", err)
		return ret, err
	}

	// Get the token from the credentials
	token, err := creds.TokenSource.Token()
	if err != nil {
		fmt.Printf("Error retrieving token: %v\n", err)
		return ret, err
	}

	// Print the access token
	fmt.Printf("Access Token: %s\n", token.AccessToken)
	return token.AccessToken, nil

}
