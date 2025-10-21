package main

import (
	"fmt"
	"log"
	"os"

	ghl "github.com/checkoutjoy/gohighlevel-go"
)

// This example shows how to use the SDK with refresh token functionality.
// Client ID and secret ARE required if you want automatic token refresh.
func main() {
	// When you need token refresh capability, provide client credentials
	client, err := ghl.NewClient(ghl.Config{
		ClientID:     os.Getenv("GHL_CLIENT_ID"),
		ClientSecret: os.Getenv("GHL_CLIENT_SECRET"),
	})
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Option 1: Use existing tokens
	accessToken := os.Getenv("GHL_ACCESS_TOKEN")
	refreshToken := os.Getenv("GHL_REFRESH_TOKEN")

	if accessToken != "" && refreshToken != "" {
		// Set both tokens with expiry (in seconds)
		client.SetTokens(accessToken, refreshToken, 3600)
		fmt.Println("Using existing access and refresh tokens")
	} else if refreshToken != "" {
		// Option 2: Refresh to get a new access token
		fmt.Println("Refreshing access token...")
		err = client.AuthorizeWithRefreshToken(refreshToken)
		if err != nil {
			log.Fatalf("Failed to refresh token: %v", err)
		}
		fmt.Printf("New access token: %s\n", client.GetAccessToken())
		fmt.Printf("New refresh token: %s\n", client.GetRefreshToken())
	} else {
		log.Fatal("Either GHL_ACCESS_TOKEN or GHL_REFRESH_TOKEN is required")
	}

	// Get location ID from environment
	locationID := os.Getenv("GHL_LOCATION_ID")
	if locationID == "" {
		log.Fatal("GHL_LOCATION_ID environment variable is required")
	}

	// Make API calls as usual
	fmt.Println("\n=== Creating Contact ===")
	contact, err := client.Contacts.Create(&ghl.CreateContactRequest{
		LocationID: locationID,
		FirstName:  "Test",
		LastName:   "User",
		Email:      "test.user@example.com",
	})
	if err != nil {
		log.Fatalf("Failed to create contact: %v", err)
	}
	fmt.Printf("Created contact: %s\n", contact.ID)

	// Clean up
	_ = client.Contacts.Delete(contact.ID)

	// If your access token expires, you can manually refresh:
	if client.GetRefreshToken() != "" {
		fmt.Println("\n=== Refreshing Token (Example) ===")
		err = client.AuthorizeWithRefreshToken(client.GetRefreshToken())
		if err != nil {
			log.Printf("Warning: Failed to refresh token: %v", err)
		} else {
			fmt.Println("Token refreshed successfully")
			// Save the new tokens for future use
			fmt.Printf("New Access Token: %s\n", client.GetAccessToken())
			fmt.Printf("New Refresh Token: %s\n", client.GetRefreshToken())
		}
	}
}
