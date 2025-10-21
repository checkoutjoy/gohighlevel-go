package gohighlevel

import (
	"fmt"
	"net/url"
	"time"
)

// ContactsService handles operations related to contacts
type ContactsService struct {
	client *Client
}

// Contact represents a GoHighLevel contact
type Contact struct {
	ID                   string             `json:"id,omitempty"`
	LocationID           string             `json:"locationId,omitempty"`
	ContactName          string             `json:"contactName,omitempty"`
	FirstName            string             `json:"firstName,omitempty"`
	LastName             string             `json:"lastName,omitempty"`
	Email                string             `json:"email,omitempty"`
	Phone                string             `json:"phone,omitempty"`
	Type                 string             `json:"type,omitempty"`
	Source               string             `json:"source,omitempty"`
	AssignedTo           string             `json:"assignedTo,omitempty"`
	Address1             string             `json:"address1,omitempty"`
	City                 string             `json:"city,omitempty"`
	State                string             `json:"state,omitempty"`
	Country              string             `json:"country,omitempty"`
	PostalCode           string             `json:"postalCode,omitempty"`
	CompanyName          string             `json:"companyName,omitempty"`
	Website              string             `json:"website,omitempty"`
	Tags                 []string           `json:"tags,omitempty"`
	DateOfBirth          string             `json:"dateOfBirth,omitempty"`
	DateAdded            time.Time          `json:"dateAdded,omitempty"`
	DateUpdated          time.Time          `json:"dateUpdated,omitempty"`
	CustomFields         []CustomField      `json:"customField,omitempty"`
	BusinessID           string             `json:"businessId,omitempty"`
	AttributionSource    *AttributionSource `json:"attributionSource,omitempty"`
	AdditionalEmails     []string           `json:"additionalEmails,omitempty"`
	AdditionalPhones     []string           `json:"additionalPhones,omitempty"`
	SSN                  string             `json:"ssn,omitempty"`
	Gender               string             `json:"gender,omitempty"`
	Timezone             string             `json:"timezone,omitempty"`
	DND                  bool               `json:"dnd,omitempty"`
	DNDSettings          *DNDSettings       `json:"dndSettings,omitempty"`
	InboundDNDSettings   *DNDSettings       `json:"inboundDndSettings,omitempty"`
	ConversationID       string             `json:"conversationId,omitempty"`
	ConversationProvider string             `json:"conversationProvider,omitempty"`
	ConversationAgencyID string             `json:"conversationAgencyId,omitempty"`
	Followers            []string           `json:"followers,omitempty"`
}

// CustomField represents a custom field on a contact
type CustomField struct {
	ID    string      `json:"id,omitempty"`
	Key   string      `json:"key,omitempty"`
	Value interface{} `json:"field_value,omitempty"`
}

// AttributionSource represents the attribution source for a contact
type AttributionSource struct {
	Campaign          string `json:"campaign,omitempty"`
	CampaignID        string `json:"campaignId,omitempty"`
	Medium            string `json:"medium,omitempty"`
	MediumID          string `json:"mediumId,omitempty"`
	Source            string `json:"source,omitempty"`
	Referrer          string `json:"referrer,omitempty"`
	AdGroup           string `json:"adGroup,omitempty"`
	AdGroupID         string `json:"adGroupId,omitempty"`
	FBCLId            string `json:"fbclid,omitempty"`
	GCLId             string `json:"gclid,omitempty"`
	MSCLKId           string `json:"msclkid,omitempty"`
	DCLID             string `json:"dclid,omitempty"`
	FBC               string `json:"fbc,omitempty"`
	FBP               string `json:"fbp,omitempty"`
	UserAgent         string `json:"userAgent,omitempty"`
	GAPClientID       string `json:"gapClientId,omitempty"`
	GoogleAnalyticsID string `json:"googleAnalyticsId,omitempty"`
}

// DNDSettings represents do not disturb settings
type DNDSettings struct {
	Call     *DNDSetting `json:"Call,omitempty"`
	SMS      *DNDSetting `json:"SMS,omitempty"`
	Email    *DNDSetting `json:"Email,omitempty"`
	WhatsApp *DNDSetting `json:"WhatsApp,omitempty"`
	GMB      *DNDSetting `json:"GMB,omitempty"`
	FB       *DNDSetting `json:"FB,omitempty"`
}

// DNDSetting represents individual DND setting
type DNDSetting struct {
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
	Code    string `json:"code,omitempty"`
}

// CreateContactRequest represents a request to create a contact
type CreateContactRequest struct {
	FirstName         string             `json:"firstName,omitempty"`
	LastName          string             `json:"lastName,omitempty"`
	Name              string             `json:"name,omitempty"`
	Email             string             `json:"email,omitempty"`
	LocationID        string             `json:"locationId"`
	Phone             string             `json:"phone,omitempty"`
	Address1          string             `json:"address1,omitempty"`
	City              string             `json:"city,omitempty"`
	State             string             `json:"state,omitempty"`
	PostalCode        string             `json:"postalCode,omitempty"`
	Country           string             `json:"country,omitempty"`
	CompanyName       string             `json:"companyName,omitempty"`
	Website           string             `json:"website,omitempty"`
	Source            string             `json:"source,omitempty"`
	Tags              []string           `json:"tags,omitempty"`
	CustomFields      []CustomField      `json:"customField,omitempty"`
	AttributionSource *AttributionSource `json:"attributionSource,omitempty"`
}

// UpdateContactRequest represents a request to update a contact
type UpdateContactRequest struct {
	FirstName         string             `json:"firstName,omitempty"`
	LastName          string             `json:"lastName,omitempty"`
	Name              string             `json:"name,omitempty"`
	Email             string             `json:"email,omitempty"`
	Phone             string             `json:"phone,omitempty"`
	Address1          string             `json:"address1,omitempty"`
	City              string             `json:"city,omitempty"`
	State             string             `json:"state,omitempty"`
	PostalCode        string             `json:"postalCode,omitempty"`
	Country           string             `json:"country,omitempty"`
	CompanyName       string             `json:"companyName,omitempty"`
	Website           string             `json:"website,omitempty"`
	Source            string             `json:"source,omitempty"`
	Tags              []string           `json:"tags,omitempty"`
	CustomFields      []CustomField      `json:"customField,omitempty"`
	AttributionSource *AttributionSource `json:"attributionSource,omitempty"`
}

// UpsertContactRequest represents a request to upsert a contact
type UpsertContactRequest struct {
	FirstName         string             `json:"firstName,omitempty"`
	LastName          string             `json:"lastName,omitempty"`
	Name              string             `json:"name,omitempty"`
	Email             string             `json:"email,omitempty"`
	LocationID        string             `json:"locationId"`
	Phone             string             `json:"phone,omitempty"`
	Address1          string             `json:"address1,omitempty"`
	City              string             `json:"city,omitempty"`
	State             string             `json:"state,omitempty"`
	PostalCode        string             `json:"postalCode,omitempty"`
	Country           string             `json:"country,omitempty"`
	CompanyName       string             `json:"companyName,omitempty"`
	Website           string             `json:"website,omitempty"`
	Source            string             `json:"source,omitempty"`
	Tags              []string           `json:"tags,omitempty"`
	CustomFields      []CustomField      `json:"customField,omitempty"`
	AttributionSource *AttributionSource `json:"attributionSource,omitempty"`
}

// GetContactsOptions represents query options for listing contacts
type GetContactsOptions struct {
	LocationID   string
	Query        string
	Limit        int
	Skip         int
	StartAfter   string
	StartAfterID string
}

// ContactResponse represents a single contact API response
type ContactResponse struct {
	Contact *Contact `json:"contact,omitempty"`
}

// ContactsResponse represents a list of contacts API response
type ContactsResponse struct {
	Contacts []Contact `json:"contacts,omitempty"`
	Total    int       `json:"total,omitempty"`
	Count    int       `json:"count,omitempty"`
}

// Create creates a new contact
// Required scope: contacts.write
func (s *ContactsService) Create(req *CreateContactRequest) (*Contact, error) {
	if req.LocationID == "" {
		return nil, fmt.Errorf("locationId is required")
	}

	var result ContactResponse
	err := s.client.doRequest("POST", "/contacts/", req, &result)
	if err != nil {
		return nil, err
	}

	return result.Contact, nil
}

// Get retrieves a contact by ID
// Required scope: contacts.readonly
func (s *ContactsService) Get(contactID string) (*Contact, error) {
	if contactID == "" {
		return nil, fmt.Errorf("contactId is required")
	}

	var result ContactResponse
	err := s.client.doRequest("GET", fmt.Sprintf("/contacts/%s", contactID), nil, &result)
	if err != nil {
		return nil, err
	}

	return result.Contact, nil
}

// Update updates an existing contact
// Required scope: contacts.write
func (s *ContactsService) Update(contactID string, req *UpdateContactRequest) (*Contact, error) {
	if contactID == "" {
		return nil, fmt.Errorf("contactId is required")
	}

	var result ContactResponse
	err := s.client.doRequest("PUT", fmt.Sprintf("/contacts/%s", contactID), req, &result)
	if err != nil {
		return nil, err
	}

	return result.Contact, nil
}

// Delete deletes a contact
// Required scope: contacts.write
func (s *ContactsService) Delete(contactID string) error {
	if contactID == "" {
		return fmt.Errorf("contactId is required")
	}

	return s.client.doRequest("DELETE", fmt.Sprintf("/contacts/%s", contactID), nil, nil)
}

// Upsert creates or updates a contact based on duplicate detection settings
// Required scope: contacts.write
func (s *ContactsService) Upsert(req *UpsertContactRequest) (*Contact, error) {
	if req.LocationID == "" {
		return nil, fmt.Errorf("locationId is required")
	}

	var result ContactResponse
	err := s.client.doRequest("POST", "/contacts/upsert", req, &result)
	if err != nil {
		return nil, err
	}

	return result.Contact, nil
}

// List retrieves a list of contacts with optional filters
// Required scope: contacts.readonly
// Note: This endpoint is deprecated, use Search instead for new implementations
func (s *ContactsService) List(opts *GetContactsOptions) (*ContactsResponse, error) {
	if opts == nil {
		opts = &GetContactsOptions{}
	}

	query := url.Values{}
	if opts.LocationID != "" {
		query.Set("locationId", opts.LocationID)
	}
	if opts.Query != "" {
		query.Set("query", opts.Query)
	}
	if opts.Limit > 0 {
		query.Set("limit", fmt.Sprintf("%d", opts.Limit))
	}
	if opts.Skip > 0 {
		query.Set("skip", fmt.Sprintf("%d", opts.Skip))
	}
	if opts.StartAfter != "" {
		query.Set("startAfter", opts.StartAfter)
	}
	if opts.StartAfterID != "" {
		query.Set("startAfterId", opts.StartAfterID)
	}

	path := "/contacts/"
	if len(query) > 0 {
		path = path + "?" + query.Encode()
	}

	var result ContactsResponse
	err := s.client.doRequest("GET", path, nil, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetByBusinessID retrieves contacts by business ID
// Required scope: contacts.readonly
func (s *ContactsService) GetByBusinessID(businessID string) (*ContactsResponse, error) {
	if businessID == "" {
		return nil, fmt.Errorf("businessId is required")
	}

	var result ContactsResponse
	err := s.client.doRequest("GET", fmt.Sprintf("/contacts/business/%s", businessID), nil, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// AddTags adds tags to a contact
// Required scope: contacts.write
func (s *ContactsService) AddTags(contactID string, tags []string) error {
	if contactID == "" {
		return fmt.Errorf("contactId is required")
	}
	if len(tags) == 0 {
		return fmt.Errorf("at least one tag is required")
	}

	req := map[string][]string{"tags": tags}
	return s.client.doRequest("POST", fmt.Sprintf("/contacts/%s/tags", contactID), req, nil)
}

// RemoveTags removes tags from a contact
// Required scope: contacts.write
func (s *ContactsService) RemoveTags(contactID string, tags []string) error {
	if contactID == "" {
		return fmt.Errorf("contactId is required")
	}
	if len(tags) == 0 {
		return fmt.Errorf("at least one tag is required")
	}

	req := map[string][]string{"tags": tags}
	return s.client.doRequest("DELETE", fmt.Sprintf("/contacts/%s/tags", contactID), req, nil)
}
