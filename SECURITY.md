# Security Guide

## OAuth Client Credentials: When and Why

### Do I Need Client ID and Secret?

**Short Answer:** Only if you're doing OAuth flows or automatic token refresh.

### Detailed Breakdown

#### ❌ You DON'T Need Them If:

1. **You already have an access token**
   ```go
   client, _ := ghl.NewClient(ghl.Config{})
   client.SetAccessToken("your-token")
   ```

2. **You manage tokens externally**
   - Your backend handles OAuth and just gives you tokens
   - You refresh tokens through a separate service
   - You have a token management system

3. **You're building client-side applications**
   - Mobile apps (iOS, Android)
   - Single-Page Applications (React, Vue, Angular)
   - Browser extensions
   - Desktop applications distributed to users

   ⚠️ **NEVER embed client secrets in client-side code!**

#### ✅ You DO Need Them If:

1. **Implementing full OAuth flow**
   ```go
   client, _ := ghl.NewClient(ghl.Config{
       ClientID:     "id",
       ClientSecret: "secret",
   })
   err := client.AuthorizeWithCode("code", "redirect")
   ```

2. **Automatic token refresh**
   ```go
   // SDK can refresh tokens automatically
   err := client.AuthorizeWithRefreshToken(refreshToken)
   ```

3. **Server-side application**
   - Backend API services
   - Server-to-server integrations
   - Webhook handlers
   - Scheduled jobs/workers

## Is It Secure to Have Client ID and Secret?

### Yes, BUT ONLY in the right context:

#### ✅ Secure Scenarios:

1. **Server-side applications**
   - Backend APIs running on your infrastructure
   - Environment variables (never hardcoded)
   - Secrets management systems (AWS Secrets Manager, HashiCorp Vault, etc.)

2. **Properly secured**
   ```bash
   # Good: Environment variables
   export GHL_CLIENT_ID="..."
   export GHL_CLIENT_SECRET="..."

   # Good: Secrets manager
   secret := secretsmanager.GetSecret("ghl-credentials")

   # BAD: Hardcoded
   clientSecret := "abc123..." // NEVER DO THIS
   ```

#### ❌ Insecure Scenarios:

1. **Client-side applications** - NEVER store secrets in:
   - Mobile app binaries
   - JavaScript bundles
   - Desktop app installers
   - Browser extensions
   - Any code that gets distributed to end users

2. **Version control**
   ```bash
   # Add to .gitignore
   .env
   credentials.json
   ```

3. **Public repositories**
   - Even in private repos, use environment variables
   - Never commit secrets, even temporarily

## Is This Normal?

**Yes, this is standard OAuth 2.0 practice:**

### The Two-Token System

1. **Access Token** (short-lived, ~1 hour)
   - Used for API calls
   - Can be exposed to frontend (with proper security)
   - Expires quickly
   - Limited scope/permissions

2. **Refresh Token** (long-lived, days/months)
   - Used to get new access tokens
   - Should be stored securely
   - Requires client credentials to use
   - Never exposed to frontend

### Why SDK Needs Client Credentials for Refresh

To refresh an access token, OAuth 2.0 requires:
```
POST /oauth/token
{
  "grant_type": "refresh_token",
  "refresh_token": "...",
  "client_id": "...",
  "client_secret": "..."  ← This is why!
}
```

The API verifies that:
1. The refresh token is valid
2. It belongs to this client (verified by client_id + client_secret)
3. The client is authorized to refresh tokens

## Recommended Architecture Patterns

### Pattern 1: Frontend + Backend (Most Secure)

```
┌─────────────┐
│   Frontend  │ ← Only has access token
│  (Browser)  │ ← No client credentials
└──────┬──────┘
       │ API calls with access token
       │
┌──────▼──────┐
│   Backend   │ ← Has client credentials
│   (Server)  │ ← Handles token refresh
└──────┬──────┘
       │ OAuth + Refresh
       │
┌──────▼──────┐
│ GoHighLevel │
│     API     │
└─────────────┘
```

**Implementation:**
```go
// Backend service
client, _ := ghl.NewClient(ghl.Config{
    ClientID:     os.Getenv("GHL_CLIENT_ID"),
    ClientSecret: os.Getenv("GHL_CLIENT_SECRET"),
})

// When frontend token expires, backend refreshes:
func refreshUserToken(userRefreshToken string) (string, error) {
    err := client.AuthorizeWithRefreshToken(userRefreshToken)
    if err != nil {
        return "", err
    }
    return client.GetAccessToken(), nil
}
```

### Pattern 2: Server-to-Server (Simple)

```
┌─────────────┐
│ Your Server │ ← Has everything
│             │ ← Client credentials + tokens
└──────┬──────┘
       │ Direct API calls
       │
┌──────▼──────┐
│ GoHighLevel │
│     API     │
└─────────────┘
```

**Implementation:**
```go
client, _ := ghl.NewClient(ghl.Config{
    ClientID:     os.Getenv("GHL_CLIENT_ID"),
    ClientSecret: os.Getenv("GHL_CLIENT_SECRET"),
})

// Full control over OAuth and refresh
client.AuthorizeWithCode(code, redirectURI)
// Or
client.SetTokens(accessToken, refreshToken, expiresIn)
```

### Pattern 3: Pre-authenticated (Simplest)

```
┌─────────────┐
│    Your     │ ← Only has access token
│ Application │ ← No OAuth, no refresh
└──────┬──────┘
       │ Simple API calls
       │
┌──────▼──────┐
│ GoHighLevel │
│     API     │
└─────────────┘
```

**Implementation:**
```go
// Simplest possible usage
client, _ := ghl.NewClient(ghl.Config{})
client.SetAccessToken(os.Getenv("GHL_ACCESS_TOKEN"))

// Just make API calls
contact, _ := client.Contacts.Get("contact-id")
```

**Use this when:**
- Tokens are managed externally
- Short-lived scripts/tools
- You manually refresh tokens
- Testing/development

## Best Practices Summary

1. ✅ **Use environment variables** for all secrets
2. ✅ **Use the simplest pattern** for your use case
3. ✅ **Keep client secrets on servers only**
4. ✅ **Use short-lived access tokens**
5. ✅ **Rotate credentials periodically**
6. ❌ **Never commit secrets to git**
7. ❌ **Never embed secrets in client apps**
8. ❌ **Never hardcode credentials**

## Example: Migration to Secure Pattern

### Before (Insecure)
```go
// BAD: Client credentials in code
client, _ := ghl.NewClient(ghl.Config{
    ClientID:     "hardcoded-id",
    ClientSecret: "hardcoded-secret",
})
```

### After (Secure)
```go
// GOOD: Only what you need
client, _ := ghl.NewClient(ghl.Config{})
client.SetAccessToken(os.Getenv("GHL_ACCESS_TOKEN"))
```

Or if you need refresh:
```go
// GOOD: Credentials from environment
client, _ := ghl.NewClient(ghl.Config{
    ClientID:     os.Getenv("GHL_CLIENT_ID"),
    ClientSecret: os.Getenv("GHL_CLIENT_SECRET"),
})
```

## Questions?

If you're unsure which pattern to use:

1. **Do you already have access tokens?** → Use Pattern 3 (simplest)
2. **Do you need automatic token refresh?** → Use Pattern 2 (server-to-server)
3. **Building a web app with frontend?** → Use Pattern 1 (frontend + backend)

The SDK now supports all three patterns securely!
