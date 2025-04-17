package pkg

import (
	"context"
	"fmt"

	"google.golang.org/api/idtoken"
)

func GetAuth(targetAudience string) (string, error) {
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
