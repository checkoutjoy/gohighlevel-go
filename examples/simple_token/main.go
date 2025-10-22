package main

import (
	"fmt"
	"log"
	"os"

	ghl "github.com/checkoutjoy/gohighlevel-go"
)

// This example shows how to use the SDK when you already have an access token.
// You don't need to provide client ID or secret for basic API operations.
func main() {
	// Get access token from environment
	accessToken := os.Getenv("GHL_ACCESS_TOKEN")
	if accessToken == "" {
		log.Fatal("GHL_ACCESS_TOKEN environment variable is required")
	}

	// Create a client with just the access token (simpler and more secure)
	client, err := ghl.NewClient(ghl.Config{
		AccessToken: accessToken,
	})
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Get location ID from environment
	locationID := os.Getenv("GHL_LOCATION_ID")
	if locationID == "" {
		log.Fatal("GHL_LOCATION_ID environment variable is required")
	}

	// Now you can make API calls
	fmt.Println("=== Creating Contact ===")
	contact, err := client.Contacts.Create(&ghl.CreateContactRequest{
		LocationID:  locationID,
		FirstName:   "Jane",
		LastName:    "Smith",
		Email:       "jane.smith@example.com",
		Phone:       "+1234567890",
		CompanyName: "Tech Corp",
		Tags:        []string{"api-user", "example"},
	})
	if err != nil {
		log.Fatalf("Failed to create contact: %v", err)
	}

	fmt.Printf("Created contact: %s (ID: %s)\n", contact.ContactName, contact.ID)

	// Get the contact
	fmt.Println("\n=== Getting Contact ===")
	retrieved, err := client.Contacts.Get(contact.ID)
	if err != nil {
		log.Fatalf("Failed to get contact: %v", err)
	}
	fmt.Printf("Retrieved: %s %s (%s)\n", retrieved.FirstName, retrieved.LastName, retrieved.Email)

	// Clean up
	fmt.Println("\n=== Cleaning Up ===")
	err = client.Contacts.Delete(contact.ID)
	if err != nil {
		log.Fatalf("Failed to delete contact: %v", err)
	}
	fmt.Println("Contact deleted successfully")

	fmt.Println("\n=== Example completed successfully ===")
}
