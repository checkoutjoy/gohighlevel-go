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

// TokenRefreshCallback is called whenever tokens are automatically refreshed due to 401 errors.
// This allows you to save the new tokens to your external storage (database, cache, etc.).
type TokenRefreshCallback func(accessToken, refreshToken string, expiresIn int)

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

	// LocationID is the default location ID for API requests
	locationID string

	// Token refresh configuration
	onTokenRefresh   TokenRefreshCallback
	autoRefreshOn401 bool

	// Resources
	Contacts *ContactsService
}

// Config holds configuration for the GoHighLevel client
type Config struct {
	ClientID         string
	ClientSecret     string
	AccessToken      string
	RefreshToken     string
	LocationID       string
	BaseURL          string
	HTTPClient       *http.Client
	OnTokenRefresh   TokenRefreshCallback // Called when tokens are automatically refreshed on 401
	AutoRefreshOn401 bool                 // Enable automatic token refresh on 401 errors (default: false)
}

// NewClient creates a new GoHighLevel API client.
// ClientID and ClientSecret are optional - only required for OAuth flows and token refresh.
// If you only need to make API calls with an existing access token, you can omit them.
func NewClient(config Config) (*Client, error) {

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
		BaseURL:          baseURL,
		HTTPClient:       httpClient,
		clientID:         config.ClientID,
		clientSecret:     config.ClientSecret,
		accessToken:      config.AccessToken,
		refreshToken:     config.RefreshToken,
		locationID:       config.LocationID,
		onTokenRefresh:   config.OnTokenRefresh,
		autoRefreshOn401: config.AutoRefreshOn401,
	}

	// Initialize services
	c.Contacts = &ContactsService{client: c}

	return c, nil
}

// AuthorizeWithCode exchanges an authorization code for an access token.
// Requires ClientID and ClientSecret to be set in the client config.
func (c *Client) AuthorizeWithCode(code, redirectURI string) error {
	if c.clientID == "" || c.clientSecret == "" {
		return fmt.Errorf("clientID and clientSecret are required for OAuth authorization")
	}
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

// AuthorizeWithRefreshToken refreshes the access token using a refresh token.
// Requires ClientID and ClientSecret to be set in the client config.
func (c *Client) AuthorizeWithRefreshToken(refreshToken string) error {
	if c.clientID == "" || c.clientSecret == "" {
		return fmt.Errorf("clientID and clientSecret are required for token refresh")
	}
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

// SetLocationID sets the default location ID for API requests
func (c *Client) SetLocationID(locationID string) {
	c.locationID = locationID
}

// GetLocationID returns the current default location ID
func (c *Client) GetLocationID() string {
	return c.locationID
}

// refreshTokenInternal is an internal method that refreshes the token and calls the callback
// This is used for automatic token refresh on 401 errors
func (c *Client) refreshTokenInternal(refreshToken string) error {
	if c.clientID == "" || c.clientSecret == "" {
		return fmt.Errorf("clientID and clientSecret are required for token refresh")
	}

	data := url.Values{}
	data.Set("client_id", c.clientID)
	data.Set("client_secret", c.clientSecret)
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)

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

	// Update tokens
	c.tokenMutex.Lock()
	c.accessToken = tokenResp.AccessToken
	c.refreshToken = tokenResp.RefreshToken
	if tokenResp.ExpiresIn > 0 {
		c.tokenExpiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	}
	c.tokenMutex.Unlock()

	// Call the callback if set (this is automatic refresh, so always call it)
	if c.onTokenRefresh != nil {
		c.onTokenRefresh(tokenResp.AccessToken, tokenResp.RefreshToken, tokenResp.ExpiresIn)
	}

	return nil
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
	// First attempt
	statusCode, respBody, err := c.executeRequest(method, path, body)

	// Check if we got a 401 and should auto-refresh
	if statusCode == http.StatusUnauthorized && c.autoRefreshOn401 {
		// Check if we have the necessary credentials to refresh
		c.tokenMutex.RLock()
		hasRefreshToken := c.refreshToken != ""
		hasCredentials := c.clientID != "" && c.clientSecret != ""
		currentRefreshToken := c.refreshToken
		c.tokenMutex.RUnlock()

		if hasRefreshToken && hasCredentials {
			// Attempt to refresh the token
			refreshErr := c.refreshTokenInternal(currentRefreshToken)
			if refreshErr != nil {
				// Refresh failed, return original error
				return fmt.Errorf("API request failed with status %d: %s (token refresh failed: %w)", statusCode, string(respBody), refreshErr)
			}

			// Retry the request with new token
			statusCode, respBody, err = c.executeRequest(method, path, body)
		}
	}

	// Handle the final response
	if err != nil {
		return err
	}

	if statusCode < 200 || statusCode >= 300 {
		return fmt.Errorf("API request failed with status %d: %s", statusCode, string(respBody))
	}

	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("failed to parse response: %w", err)
		}
	}

	return nil
}

// executeRequest performs the actual HTTP request and returns status code, body, and error
func (c *Client) executeRequest(method, path string, body interface{}) (int, []byte, error) {
	c.tokenMutex.RLock()
	token := c.accessToken
	c.tokenMutex.RUnlock()

	if token == "" {
		return 0, nil, fmt.Errorf("no access token available, please authorize first")
	}

	var bodyReader io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return 0, nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewBuffer(jsonData)
	}

	fullURL := c.BaseURL + path
	req, err := http.NewRequest(method, fullURL, bodyReader)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return 0, nil, fmt.Errorf("request failed: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, nil, fmt.Errorf("failed to read response: %w", err)
	}

	return resp.StatusCode, respBody, nil
}
