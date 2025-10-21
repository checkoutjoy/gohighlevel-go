package main

import (
	"fmt"
	"log"
	"os"

	ghl "github.com/checkoutjoy/gohighlevel-go"
)

func main() {
	// Initialize the client with OAuth credentials
	client, err := ghl.NewClient(ghl.Config{
		ClientID:     os.Getenv("GHL_CLIENT_ID"),
		ClientSecret: os.Getenv("GHL_CLIENT_SECRET"),
	})
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Option 1: Use an existing access token
	accessToken := os.Getenv("GHL_ACCESS_TOKEN")
	if accessToken != "" {
		client.SetAccessToken(accessToken)
	} else {
		// Option 2: Authorize with OAuth code
		authCode := os.Getenv("GHL_AUTH_CODE")
		redirectURI := os.Getenv("GHL_REDIRECT_URI")

		err = client.AuthorizeWithCode(authCode, redirectURI)
		if err != nil {
			log.Fatalf("Failed to authorize: %v", err)
		}

		fmt.Printf("Access Token: %s\n", client.GetAccessToken())
		fmt.Printf("Refresh Token: %s\n", client.GetRefreshToken())
	}

	// Get location ID from environment
	locationID := os.Getenv("GHL_LOCATION_ID")

	// Create a new contact
	fmt.Println("\n=== Creating Contact ===")
	contact, err := client.Contacts.Create(&ghl.CreateContactRequest{
		LocationID:  locationID,
		FirstName:   "John",
		LastName:    "Doe",
		Email:       "john.doe@example.com",
		Phone:       "+1234567890",
		CompanyName: "Acme Corp",
		Tags:        []string{"lead", "website"},
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

	// Update the contact
	fmt.Println("\n=== Updating Contact ===")
	updated, err := client.Contacts.Update(contact.ID, &ghl.UpdateContactRequest{
		CompanyName: "Acme Corporation Inc.",
		Tags:        []string{"lead", "website", "qualified"},
	})
	if err != nil {
		log.Fatalf("Failed to update contact: %v", err)
	}
	fmt.Printf("Updated company to: %s\n", updated.CompanyName)

	// Add tags
	fmt.Println("\n=== Adding Tags ===")
	err = client.Contacts.AddTags(contact.ID, []string{"high-priority", "demo-requested"})
	if err != nil {
		log.Fatalf("Failed to add tags: %v", err)
	}
	fmt.Println("Tags added successfully")

	// List contacts
	fmt.Println("\n=== Listing Contacts ===")
	contacts, err := client.Contacts.List(&ghl.GetContactsOptions{
		LocationID: locationID,
		Limit:      10,
	})
	if err != nil {
		log.Fatalf("Failed to list contacts: %v", err)
	}
	fmt.Printf("Found %d contacts (Total: %d)\n", contacts.Count, contacts.Total)

	// Upsert a contact
	fmt.Println("\n=== Upserting Contact ===")
	upserted, err := client.Contacts.Upsert(&ghl.UpsertContactRequest{
		LocationID: locationID,
		Email:      "jane.smith@example.com",
		FirstName:  "Jane",
		LastName:   "Smith",
		Phone:      "+1987654321",
	})
	if err != nil {
		log.Fatalf("Failed to upsert contact: %v", err)
	}
	fmt.Printf("Upserted contact: %s (ID: %s)\n", upserted.ContactName, upserted.ID)

	// Remove tags
	fmt.Println("\n=== Removing Tags ===")
	err = client.Contacts.RemoveTags(contact.ID, []string{"lead"})
	if err != nil {
		log.Fatalf("Failed to remove tags: %v", err)
	}
	fmt.Println("Tags removed successfully")

	// Clean up - delete the created contact
	fmt.Println("\n=== Deleting Contact ===")
	err = client.Contacts.Delete(contact.ID)
	if err != nil {
		log.Fatalf("Failed to delete contact: %v", err)
	}
	fmt.Println("Contact deleted successfully")

	// Delete the upserted contact too
	err = client.Contacts.Delete(upserted.ID)
	if err != nil {
		log.Fatalf("Failed to delete upserted contact: %v", err)
	}

	fmt.Println("\n=== Example completed successfully ===")
}
