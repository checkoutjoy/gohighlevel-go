# GoHighLevel Go SDK

A comprehensive, production-ready Go SDK for the [GoHighLevel API](https://marketplace.gohighlevel.com/). This SDK provides a clean, resource-based interface for interacting with GoHighLevel's CRM platform, with full OAuth 2.0 support.

**Built by [CheckoutJoy](https://checkoutjoy.com) - Simplifying payment solutions**

[![CI](https://github.com/checkoutjoy/gohighlevel-go/actions/workflows/ci.yml/badge.svg)](https://github.com/checkoutjoy/gohighlevel-go/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/checkoutjoy/gohighlevel-go)](https://goreportcard.com/report/github.com/checkoutjoy/gohighlevel-go)
[![GoDoc](https://godoc.org/github.com/checkoutjoy/gohighlevel-go?status.svg)](https://godoc.org/github.com/checkoutjoy/gohighlevel-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Features

- **OAuth 2.0 Authentication** - Full support for GoHighLevel's OAuth flow
- **Resource-Based API** - Clean, intuitive interface organized by resources
- **Type-Safe** - Comprehensive Go types for all API entities
- **Connection Pooling** - Efficient HTTP client with connection reuse
- **Comprehensive Testing** - Full integration test suite
- **Production Ready** - Built with best practices for Go SDKs

## Installation

```bash
go get github.com/checkoutjoy/gohighlevel-go
```

## Quick Start

### Simple Usage (With Access Token Only)

If you already have an access token, you don't need to provide client credentials:

```go
package main

import (
    "fmt"
    "log"

    ghl "github.com/checkoutjoy/gohighlevel-go"
)

func main() {
    // Create a client without OAuth credentials (simpler and more secure)
    client, err := ghl.NewClient(ghl.Config{})
    if err != nil {
        log.Fatal(err)
    }

    // Just set your access token
    client.SetAccessToken("your-access-token")

    // Create a contact
    contact, err := client.Contacts.Create(&ghl.CreateContactRequest{
        LocationID: "location-id",
        FirstName:  "John",
        LastName:   "Doe",
        Email:      "john.doe@example.com",
        Phone:      "+1234567890",
        Tags:       []string{"lead", "website"},
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Created contact: %s\n", contact.ID)
}
```

## Authentication

The SDK supports OAuth 2.0 authentication with multiple authorization methods.

### Important: When Do You Need Client ID & Secret?

**You DON'T need them if:**
- You already have a valid access token
- You manage token refresh externally
- You're building a client-side application (they should NEVER be exposed client-side)

**You DO need them if:**
- You're implementing OAuth authorization code flow
- You want the SDK to automatically refresh tokens using a refresh token
- You're building a server-side application that handles OAuth

### Security Best Practices

⚠️ **Never expose client secrets in client-side code** (mobile apps, SPAs, browser extensions)
✅ **Only use client secrets in secure server-side environments**
✅ **If you only have access tokens, you can use the SDK without any client credentials**

### Method 1: Simple - Using Just an Access Token (Recommended)

Most secure for applications that already have tokens:

```go
// No client credentials needed!
client, _ := ghl.NewClient(ghl.Config{})
client.SetAccessToken("your-access-token")

// Start making API calls immediately
contact, err := client.Contacts.Get("contact-id")
```

### Method 2: With Token Refresh Capability

If you need the SDK to refresh tokens automatically:

```go
// Provide client credentials for token refresh
client, _ := ghl.NewClient(ghl.Config{
    ClientID:     "your-client-id",
    ClientSecret: "your-client-secret",
})

// Set initial tokens
client.SetTokens("access-token", "refresh-token", 3600)

// Later, refresh when needed
err := client.AuthorizeWithRefreshToken(client.GetRefreshToken())
```

### Method 3: Full OAuth Authorization Code Flow

For server-side apps implementing OAuth from scratch:

```go
client, _ := ghl.NewClient(ghl.Config{
    ClientID:     "your-client-id",
    ClientSecret: "your-client-secret",
})

// Exchange authorization code for access token
err := client.AuthorizeWithCode("auth-code", "redirect-uri")
```

### Method 4: Manual Token Management

```go
// Set tokens with expiry
client.SetTokens(
    "access-token",
    "refresh-token",
    3600, // expires in seconds
)

// Retrieve current tokens
accessToken := client.GetAccessToken()
refreshToken := client.GetRefreshToken()
```

## Resources

### Contacts

The Contacts resource provides full CRUD operations for managing contacts in GoHighLevel.

#### Create a Contact

```go
contact, err := client.Contacts.Create(&ghl.CreateContactRequest{
    LocationID:  "location-id",
    FirstName:   "Jane",
    LastName:    "Smith",
    Email:       "jane@example.com",
    Phone:       "+1987654321",
    CompanyName: "Acme Corp",
    Tags:        []string{"customer", "vip"},
    CustomFields: []ghl.CustomField{
        {Key: "industry", Value: "Technology"},
    },
})
```

**Required Scope:** `contacts.write`

#### Get a Contact

```go
contact, err := client.Contacts.Get("contact-id")
```

**Required Scope:** `contacts.readonly`

#### Update a Contact

```go
updated, err := client.Contacts.Update("contact-id", &ghl.UpdateContactRequest{
    LastName:    "Johnson",
    CompanyName: "New Company Inc",
})
```

**Required Scope:** `contacts.write`

#### Delete a Contact

```go
err := client.Contacts.Delete("contact-id")
```

**Required Scope:** `contacts.write`

#### Upsert a Contact

Create or update a contact based on duplicate detection settings:

```go
contact, err := client.Contacts.Upsert(&ghl.UpsertContactRequest{
    LocationID: "location-id",
    Email:      "user@example.com",
    FirstName:  "John",
    LastName:   "Doe",
})
```

**Required Scope:** `contacts.write`

**Note:** The Upsert API adheres to the "Allow Duplicate Contact" setting at the Location level. If both email and phone match different contacts, it updates the first field in the configured sequence.

#### List Contacts

```go
contacts, err := client.Contacts.List(&ghl.GetContactsOptions{
    LocationID: "location-id",
    Limit:      50,
    Query:      "search-term",
})

fmt.Printf("Found %d of %d contacts\n", contacts.Count, contacts.Total)
for _, contact := range contacts.Contacts {
    fmt.Printf("- %s (%s)\n", contact.ContactName, contact.Email)
}
```

**Required Scope:** `contacts.readonly`

**Note:** This endpoint is deprecated. Use the Search Contacts endpoint for new implementations.

#### Get Contacts by Business ID

```go
contacts, err := client.Contacts.GetByBusinessID("business-id")
```

**Required Scope:** `contacts.readonly`

### Contact Tags

#### Add Tags to a Contact

```go
err := client.Contacts.AddTags("contact-id", []string{"qualified", "hot-lead"})
```

**Required Scope:** `contacts.write`

#### Remove Tags from a Contact

```go
err := client.Contacts.RemoveTags("contact-id", []string{"cold-lead"})
```

**Required Scope:** `contacts.write`


## OAuth Scopes

The following OAuth scopes are required for different operations:

| Scope | Description | Operations |
|-------|-------------|------------|
| `contacts.readonly` | Read access to contacts | Get Contact, List Contacts, Get Contacts by Business ID |
| `contacts.write` | Write access to contacts | Create Contact, Update Contact, Delete Contact, Upsert Contact, Add Tags, Remove Tags |

### Requesting Scopes

When setting up your OAuth application, request the appropriate scopes based on your needs:

- **Read-only access:** `contacts.readonly`
- **Full access:** `contacts.readonly contacts.write`

## Configuration

### Custom HTTP Client

Provide your own HTTP client with custom settings:

```go
import (
    "net/http"
    "time"
)

customClient := &http.Client{
    Timeout: 60 * time.Second,
    Transport: &http.Transport{
        MaxIdleConns:        200,
        MaxIdleConnsPerHost: 20,
        IdleConnTimeout:     120 * time.Second,
    },
}

client, err := ghl.NewClient(ghl.Config{
    ClientID:     "your-client-id",
    ClientSecret: "your-client-secret",
    HTTPClient:   customClient,
})
```

### Custom Base URL

For testing or custom deployments:

```go
client, err := ghl.NewClient(ghl.Config{
    ClientID:     "your-client-id",
    ClientSecret: "your-client-secret",
    BaseURL:      "https://custom-api.example.com",
})
```

## Development

### Prerequisites

- Go 1.24 or higher
- Make (optional, for using Makefile)

### Building

```bash
make build
```

Or without Make:

```bash
go build ./...
```

### Running Tests

#### Unit Tests

```bash
make test-unit
```

Or:

```bash
go test -short ./...
```

#### Integration Tests

Integration tests require valid GoHighLevel credentials:

```bash
export GHL_CLIENT_ID="your-client-id"
export GHL_CLIENT_SECRET="your-client-secret"
export GHL_ACCESS_TOKEN="your-access-token"
export GHL_LOCATION_ID="your-location-id"

make test-integration
```

Or:

```bash
go test -v ./...
```

#### All Tests with Coverage

```bash
make test
```

This generates a coverage report in `coverage.html`.

### Linting

```bash
make lint
```

Or:

```bash
golangci-lint run ./...
```

### Formatting

```bash
make fmt
```

## Examples

See the [examples](./examples) directory for complete working examples:

### 1. Simple Usage (Access Token Only)
[examples/simple_token](./examples/simple_token) - Use the SDK with just an access token (no client credentials needed)

```bash
export GHL_ACCESS_TOKEN="your-access-token"
export GHL_LOCATION_ID="your-location-id"
go run examples/simple_token/main.go
```

### 2. With Token Refresh
[examples/refresh_token](./examples/refresh_token) - Use the SDK with automatic token refresh capability

```bash
export GHL_CLIENT_ID="your-client-id"
export GHL_CLIENT_SECRET="your-client-secret"
export GHL_ACCESS_TOKEN="your-access-token"
export GHL_REFRESH_TOKEN="your-refresh-token"
export GHL_LOCATION_ID="your-location-id"
go run examples/refresh_token/main.go
```

### 3. Full OAuth Flow
[examples/basic](./examples/basic) - Comprehensive example covering all major operations including OAuth flow

```bash
export GHL_CLIENT_ID="your-client-id"
export GHL_CLIENT_SECRET="your-client-secret"
export GHL_ACCESS_TOKEN="your-access-token"
export GHL_LOCATION_ID="your-location-id"
go run examples/basic/main.go
```

## API Documentation

For detailed API documentation, refer to:

- [GoHighLevel API Documentation](https://marketplace.gohighlevel.com/docs)
- [OAuth 2.0 Guide](https://marketplace.gohighlevel.com/docs/ghl/oauth/o-auth-2-0)
- [Contacts API](https://marketplace.gohighlevel.com/docs/ghl/contacts/contacts)
- [Contact Tags API](https://marketplace.gohighlevel.com/docs/ghl/contacts/tags)

## Error Handling

The SDK returns descriptive errors for all operations:

```go
contact, err := client.Contacts.Get("invalid-id")
if err != nil {
    // Handle error
    log.Printf("Failed to get contact: %v", err)
    return
}
```

Common error scenarios:
- Missing or invalid access token
- Invalid request parameters
- API rate limiting
- Network errors
- Invalid credentials

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

For issues and questions:

- GitHub Issues: [https://github.com/checkoutjoy/gohighlevel-go/issues](https://github.com/checkoutjoy/gohighlevel-go/issues)
- GoHighLevel API Support: [https://marketplace.gohighlevel.com/](https://marketplace.gohighlevel.com/)

Built by [CheckoutJoy](https://checkoutjoy.com)

