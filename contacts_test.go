package gohighlevel

import (
	"os"
	"testing"
	"time"
)

// Integration tests for Contacts API
// These tests require valid credentials set via environment variables:
// - GHL_CLIENT_ID
// - GHL_CLIENT_SECRET
// - GHL_ACCESS_TOKEN (or GHL_AUTH_CODE and GHL_REDIRECT_URI)
// - GHL_LOCATION_ID

func setupTestClient(t *testing.T) *Client {
	clientID := os.Getenv("GHL_CLIENT_ID")
	clientSecret := os.Getenv("GHL_CLIENT_SECRET")
	accessToken := os.Getenv("GHL_ACCESS_TOKEN")

	if clientID == "" || clientSecret == "" {
		t.Skip("Skipping integration tests: GHL_CLIENT_ID and GHL_CLIENT_SECRET not set")
	}

	client, err := NewClient(Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Try to authenticate if we have credentials
	if accessToken != "" {
		client.SetAccessToken(accessToken)
	} else {
		authCode := os.Getenv("GHL_AUTH_CODE")
		redirectURI := os.Getenv("GHL_REDIRECT_URI")

		if authCode != "" {
			err = client.AuthorizeWithCode(authCode, redirectURI)
			if err != nil {
				t.Fatalf("Failed to authorize with code: %v", err)
			}
		} else {
			t.Skip("Skipping integration tests: No access token or auth code provided")
		}
	}

	return client
}

func getTestLocationID(t *testing.T) string {
	locationID := os.Getenv("GHL_LOCATION_ID")
	if locationID == "" {
		t.Skip("Skipping test: GHL_LOCATION_ID not set")
	}
	return locationID
}

func TestContactsIntegration_Create(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := setupTestClient(t)
	locationID := getTestLocationID(t)

	// Create a test contact
	req := &CreateContactRequest{
		LocationID:  locationID,
		FirstName:   "Test",
		LastName:    "Contact",
		Email:       "test+" + time.Now().Format("20060102150405") + "@example.com",
		Phone:       "+15555551234",
		CompanyName: "Test Company",
		Tags:        []string{"test", "integration"},
	}

	contact, err := client.Contacts.Create(req)
	if err != nil {
		t.Fatalf("Failed to create contact: %v", err)
	}

	if contact.ID == "" {
		t.Error("Created contact has no ID")
	}

	if contact.FirstName != req.FirstName {
		t.Errorf("Expected first name %s, got %s", req.FirstName, contact.FirstName)
	}

	t.Logf("Created contact with ID: %s", contact.ID)

	// Clean up - delete the test contact
	defer func() {
		err := client.Contacts.Delete(contact.ID)
		if err != nil {
			t.Logf("Warning: Failed to delete test contact: %v", err)
		}
	}()
}

func TestContactsIntegration_Get(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := setupTestClient(t)
	locationID := getTestLocationID(t)

	// First create a contact
	createReq := &CreateContactRequest{
		LocationID: locationID,
		FirstName:  "TestGet",
		LastName:   "Contact",
		Email:      "testget+" + time.Now().Format("20060102150405") + "@example.com",
	}

	created, err := client.Contacts.Create(createReq)
	if err != nil {
		t.Fatalf("Failed to create contact: %v", err)
	}

	defer func() {
		_ = client.Contacts.Delete(created.ID)
	}()

	// Now get the contact
	contact, err := client.Contacts.Get(created.ID)
	if err != nil {
		t.Fatalf("Failed to get contact: %v", err)
	}

	if contact.ID != created.ID {
		t.Errorf("Expected contact ID %s, got %s", created.ID, contact.ID)
	}

	if contact.Email != createReq.Email {
		t.Errorf("Expected email %s, got %s", createReq.Email, contact.Email)
	}
}

func TestContactsIntegration_Update(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := setupTestClient(t)
	locationID := getTestLocationID(t)

	// Create a contact
	createReq := &CreateContactRequest{
		LocationID: locationID,
		FirstName:  "TestUpdate",
		LastName:   "Original",
		Email:      "testupdate+" + time.Now().Format("20060102150405") + "@example.com",
	}

	created, err := client.Contacts.Create(createReq)
	if err != nil {
		t.Fatalf("Failed to create contact: %v", err)
	}

	defer func() {
		_ = client.Contacts.Delete(created.ID)
	}()

	// Update the contact
	updateReq := &UpdateContactRequest{
		LastName:    "Updated",
		CompanyName: "Updated Company",
	}

	updated, err := client.Contacts.Update(created.ID, updateReq)
	if err != nil {
		t.Fatalf("Failed to update contact: %v", err)
	}

	if updated.LastName != "Updated" {
		t.Errorf("Expected last name 'Updated', got %s", updated.LastName)
	}

	if updated.CompanyName != "Updated Company" {
		t.Errorf("Expected company 'Updated Company', got %s", updated.CompanyName)
	}
}

func TestContactsIntegration_Upsert(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := setupTestClient(t)
	locationID := getTestLocationID(t)

	email := "testupsert+" + time.Now().Format("20060102150405") + "@example.com"

	// First upsert (create)
	req := &UpsertContactRequest{
		LocationID: locationID,
		Email:      email,
		FirstName:  "Upsert",
		LastName:   "Test",
	}

	contact1, err := client.Contacts.Upsert(req)
	if err != nil {
		t.Fatalf("Failed to upsert contact (create): %v", err)
	}

	defer func() {
		_ = client.Contacts.Delete(contact1.ID)
	}()

	// Second upsert (update)
	req.LastName = "Updated"
	contact2, err := client.Contacts.Upsert(req)
	if err != nil {
		t.Fatalf("Failed to upsert contact (update): %v", err)
	}

	// Should be the same contact ID
	if contact1.ID != contact2.ID {
		t.Errorf("Expected same contact ID, got %s and %s", contact1.ID, contact2.ID)
	}

	if contact2.LastName != "Updated" {
		t.Errorf("Expected last name 'Updated', got %s", contact2.LastName)
	}
}

func TestContactsIntegration_Delete(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := setupTestClient(t)
	locationID := getTestLocationID(t)

	// Create a contact to delete
	createReq := &CreateContactRequest{
		LocationID: locationID,
		FirstName:  "TestDelete",
		Email:      "testdelete+" + time.Now().Format("20060102150405") + "@example.com",
	}

	created, err := client.Contacts.Create(createReq)
	if err != nil {
		t.Fatalf("Failed to create contact: %v", err)
	}

	// Delete it
	err = client.Contacts.Delete(created.ID)
	if err != nil {
		t.Fatalf("Failed to delete contact: %v", err)
	}

	// Try to get it - should fail
	_, err = client.Contacts.Get(created.ID)
	if err == nil {
		t.Error("Expected error when getting deleted contact, got nil")
	}
}

func TestContactsIntegration_List(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := setupTestClient(t)
	locationID := getTestLocationID(t)

	// List contacts
	opts := &GetContactsOptions{
		LocationID: locationID,
		Limit:      10,
	}

	result, err := client.Contacts.List(opts)
	if err != nil {
		t.Fatalf("Failed to list contacts: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	t.Logf("Found %d contacts (total: %d)", result.Count, result.Total)
}

func TestContactsIntegration_AddTags(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := setupTestClient(t)
	locationID := getTestLocationID(t)

	// Create a contact
	createReq := &CreateContactRequest{
		LocationID: locationID,
		FirstName:  "TestTags",
		Email:      "testtags+" + time.Now().Format("20060102150405") + "@example.com",
	}

	created, err := client.Contacts.Create(createReq)
	if err != nil {
		t.Fatalf("Failed to create contact: %v", err)
	}

	defer func() {
		_ = client.Contacts.Delete(created.ID)
	}()

	// Add tags
	tags := []string{"integration-test", "automated"}
	err = client.Contacts.AddTags(created.ID, tags)
	if err != nil {
		t.Fatalf("Failed to add tags: %v", err)
	}

	t.Logf("Successfully added tags to contact %s", created.ID)
}

func TestContactsIntegration_RemoveTags(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := setupTestClient(t)
	locationID := getTestLocationID(t)

	// Create a contact with tags
	createReq := &CreateContactRequest{
		LocationID: locationID,
		FirstName:  "TestRemoveTags",
		Email:      "testremove+" + time.Now().Format("20060102150405") + "@example.com",
		Tags:       []string{"tag1", "tag2", "tag3"},
	}

	created, err := client.Contacts.Create(createReq)
	if err != nil {
		t.Fatalf("Failed to create contact: %v", err)
	}

	defer func() {
		_ = client.Contacts.Delete(created.ID)
	}()

	// Remove some tags
	tagsToRemove := []string{"tag1", "tag2"}
	err = client.Contacts.RemoveTags(created.ID, tagsToRemove)
	if err != nil {
		t.Fatalf("Failed to remove tags: %v", err)
	}

	t.Logf("Successfully removed tags from contact %s", created.ID)
}

func TestContactsIntegration_FullWorkflow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := setupTestClient(t)
	locationID := getTestLocationID(t)

	timestamp := time.Now().Format("20060102150405")

	// 1. Create a contact
	t.Log("Step 1: Creating contact")
	createReq := &CreateContactRequest{
		LocationID:  locationID,
		FirstName:   "Integration",
		LastName:    "Test",
		Email:       "integration+" + timestamp + "@example.com",
		Phone:       "+15555550001",
		CompanyName: "Test Co",
		Tags:        []string{"new-lead"},
	}

	contact, err := client.Contacts.Create(createReq)
	if err != nil {
		t.Fatalf("Failed to create contact: %v", err)
	}
	t.Logf("Created contact: %s", contact.ID)

	defer func() {
		_ = client.Contacts.Delete(contact.ID)
	}()

	// 2. Get the contact
	t.Log("Step 2: Retrieving contact")
	retrieved, err := client.Contacts.Get(contact.ID)
	if err != nil {
		t.Fatalf("Failed to get contact: %v", err)
	}
	if retrieved.Email != createReq.Email {
		t.Errorf("Email mismatch: expected %s, got %s", createReq.Email, retrieved.Email)
	}

	// 3. Update the contact
	t.Log("Step 3: Updating contact")
	updateReq := &UpdateContactRequest{
		LastName:    "Updated",
		CompanyName: "Updated Test Co",
	}
	updated, err := client.Contacts.Update(contact.ID, updateReq)
	if err != nil {
		t.Fatalf("Failed to update contact: %v", err)
	}
	if updated.LastName != "Updated" {
		t.Errorf("Last name not updated: expected 'Updated', got %s", updated.LastName)
	}

	// 4. Add tags
	t.Log("Step 4: Adding tags")
	err = client.Contacts.AddTags(contact.ID, []string{"qualified", "high-priority"})
	if err != nil {
		t.Fatalf("Failed to add tags: %v", err)
	}

	// 5. Remove tags
	t.Log("Step 5: Removing tags")
	err = client.Contacts.RemoveTags(contact.ID, []string{"new-lead"})
	if err != nil {
		t.Fatalf("Failed to remove tags: %v", err)
	}

	// 6. Delete the contact
	t.Log("Step 6: Deleting contact")
	err = client.Contacts.Delete(contact.ID)
	if err != nil {
		t.Fatalf("Failed to delete contact: %v", err)
	}

	t.Log("Full workflow completed successfully")
}
