// Package gohighlevel provides a Go SDK for the GoHighLevel API with OAuth 2.0 authentication.
package gohighlevel

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const (
	// DefaultBaseURL is the default base URL for the GoHighLevel API
	DefaultBaseURL = "https://services.leadconnectorhq.com"
	// OAuthTokenURL is the OAuth token endpoint
	OAuthTokenURL = "https://services.leadconnectorhq.com/oauth/token"
	// DefaultTimeout is the default HTTP client timeout
	DefaultTimeout = 30 * time.Second
)

// Client is the main GoHighLevel API client
type Client struct {
	// BaseURL is the base URL for API requests
	BaseURL string

	// HTTPClient is the underlying HTTP client used for requests
	HTTPClient *http.Client

	// OAuth credentials
	clientID     string
	clientSecret string

	// Access token management
	accessToken  string
	refreshToken string
	tokenExpiry  time.Time
	tokenMutex   sync.RWMutex

	// Resources
	Contacts *ContactsService
}

// Config holds configuration for the GoHighLevel client
type Config struct {
	ClientID     string
	ClientSecret string
	BaseURL      string
	HTTPClient   *http.Client
}

// NewClient creates a new GoHighLevel API client
func NewClient(config Config) (*Client, error) {
	if config.ClientID == "" {
		return nil, fmt.Errorf("clientID is required")
	}
	if config.ClientSecret == "" {
		return nil, fmt.Errorf("clientSecret is required")
	}

	baseURL := config.BaseURL
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}

	httpClient := config.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: DefaultTimeout,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
			},
		}
	}

	c := &Client{
		BaseURL:      baseURL,
		HTTPClient:   httpClient,
		clientID:     config.ClientID,
		clientSecret: config.ClientSecret,
	}

	// Initialize services
	c.Contacts = &ContactsService{client: c}

	return c, nil
}

// AuthorizeWithCode exchanges an authorization code for an access token
func (c *Client) AuthorizeWithCode(code, redirectURI string) error {
	data := url.Values{}
	data.Set("client_id", c.clientID)
	data.Set("client_secret", c.clientSecret)
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	if redirectURI != "" {
		data.Set("redirect_uri", redirectURI)
	}

	return c.fetchToken(data)
}

// AuthorizeWithRefreshToken refreshes the access token using a refresh token
func (c *Client) AuthorizeWithRefreshToken(refreshToken string) error {
	data := url.Values{}
	data.Set("client_id", c.clientID)
	data.Set("client_secret", c.clientSecret)
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)

	return c.fetchToken(data)
}

// SetAccessToken manually sets the access token
func (c *Client) SetAccessToken(token string) {
	c.tokenMutex.Lock()
	defer c.tokenMutex.Unlock()
	c.accessToken = token
}

// SetTokens manually sets both access and refresh tokens
func (c *Client) SetTokens(accessToken, refreshToken string, expiresIn int) {
	c.tokenMutex.Lock()
	defer c.tokenMutex.Unlock()
	c.accessToken = accessToken
	c.refreshToken = refreshToken
	if expiresIn > 0 {
		c.tokenExpiry = time.Now().Add(time.Duration(expiresIn) * time.Second)
	}
}

// GetAccessToken returns the current access token
func (c *Client) GetAccessToken() string {
	c.tokenMutex.RLock()
	defer c.tokenMutex.RUnlock()
	return c.accessToken
}

// GetRefreshToken returns the current refresh token
func (c *Client) GetRefreshToken() string {
	c.tokenMutex.RLock()
	defer c.tokenMutex.RUnlock()
	return c.refreshToken
}

// fetchToken fetches an access token from the OAuth endpoint
func (c *Client) fetchToken(data url.Values) error {
	req, err := http.NewRequest("POST", OAuthTokenURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create token request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to fetch token: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read token response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("token request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		TokenType    string `json:"token_type"`
		ExpiresIn    int    `json:"expires_in"`
		Scope        string `json:"scope"`
	}

	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return fmt.Errorf("failed to parse token response: %w", err)
	}

	c.tokenMutex.Lock()
	c.accessToken = tokenResp.AccessToken
	c.refreshToken = tokenResp.RefreshToken
	if tokenResp.ExpiresIn > 0 {
		c.tokenExpiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	}
	c.tokenMutex.Unlock()

	return nil
}

// doRequest performs an HTTP request with the access token
func (c *Client) doRequest(method, path string, body interface{}, result interface{}) error {
	c.tokenMutex.RLock()
	token := c.accessToken
	c.tokenMutex.RUnlock()

	if token == "" {
		return fmt.Errorf("no access token available, please authorize first")
	}

	var bodyReader io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewBuffer(jsonData)
	}

	fullURL := c.BaseURL + path
	req, err := http.NewRequest(method, fullURL, bodyReader)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("failed to parse response: %w", err)
		}
	}

	return nil
}
