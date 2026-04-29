package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mfb27/luban/internal/model"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"gorm.io/gorm"
)

// Default timeout for GitHub API requests
const githubAPITimeout = 30 * time.Second

// GitHubConfig wraps OAuth2 config for GitHub
type GitHubOAuthConfig struct {
	OAuth2Config *oauth2.Config
	HTTPClient   *http.Client // Custom HTTP client with timeout and optional proxy
}

// GitHubUserInfo represents GitHub user data from API
type GitHubUserInfo struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

// NewGitHubOAuthConfig creates a GitHub OAuth configuration with optional proxy
func NewGitHubOAuthConfig(clientID, clientSecret, redirectURL string) *GitHubOAuthConfig {
	return NewGitHubOAuthConfigWithProxy(clientID, clientSecret, redirectURL, "")
}

// NewGitHubOAuthConfigWithProxy creates a GitHub OAuth configuration with optional proxy support
func NewGitHubOAuthConfigWithProxy(clientID, clientSecret, redirectURL, proxyURL string) *GitHubOAuthConfig {
	cfg := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       []string{"user:email"},
		Endpoint:     github.Endpoint,
	}

	// Create custom HTTP client with timeout and optional proxy
	httpClient := createHTTPClient(proxyURL)

	return &GitHubOAuthConfig{
		OAuth2Config: cfg,
		HTTPClient:   httpClient,
	}
}

// createHTTPClient creates an HTTP client with timeout and optional proxy
func createHTTPClient(proxyURL string) *http.Client {
	transport := &http.Transport{
		ResponseHeaderTimeout: githubAPITimeout,
	}

	// Configure proxy if provided
	if proxyURL != "" {
		proxyURLParsed, err := url.Parse(proxyURL)
		if err == nil {
			transport.Proxy = http.ProxyURL(proxyURLParsed)
		}
	}

	return &http.Client{
		Timeout:   githubAPITimeout,
		Transport: transport,
	}
}

// GetAuthURL generates the GitHub OAuth authorization URL with state
func (c *GitHubOAuthConfig) GetAuthURL(state string) string {
	return c.OAuth2Config.AuthCodeURL(state, oauth2.AccessTypeOnline)
}

// ExchangeCode exchanges the OAuth code for an access token
func (c *GitHubOAuthConfig) ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error) {
	// Use custom HTTP client for the exchange
	ctx = context.WithValue(ctx, oauth2.HTTPClient, c.HTTPClient)
	return c.OAuth2Config.Exchange(ctx, code)
}

// GetUserInfo fetches user information from GitHub API using the access token
func (c *GitHubOAuthConfig) GetUserInfo(ctx context.Context, token *oauth2.Token) (*GitHubUserInfo, error) {
	// Create a client with the token, but use our custom HTTP client as base
	client := c.HTTPClient

	// Get user profile with Authorization header
	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.github.com/user", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitHub API error: %s - %s", resp.Status, string(body))
	}

	var userInfo GitHubUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	// If email is not public, fetch emails endpoint
	if userInfo.Email == "" {
		emailReq, err := http.NewRequestWithContext(ctx, "GET", "https://api.github.com/user/emails", nil)
		if err == nil {
			emailReq.Header.Set("Authorization", "Bearer "+token.AccessToken)
			emailResp, err := client.Do(emailReq)
			if err == nil && emailResp.StatusCode == http.StatusOK {
				defer emailResp.Body.Close()
				var emails []struct {
					Email   string `json:"email"`
					Primary bool   `json:"primary"`
				}
				if json.NewDecoder(emailResp.Body).Decode(&emails) == nil {
					for _, e := range emails {
						if e.Primary {
							userInfo.Email = e.Email
							break
						}
					}
					// Fallback to first email if no primary
					if userInfo.Email == "" && len(emails) > 0 {
						userInfo.Email = emails[0].Email
					}
				}
			}
		}
	}

	return &userInfo, nil
}

// GenerateState generates a random state string for OAuth CSRF protection
func GenerateState() string {
	return uuid.New().String()
}

// GitHubLoginRequest is the request for direct GitHub login (SPA flow)
type GitHubLoginRequest struct {
	Code string `json:"code" binding:"required"`
}

// HandleGitHubLogin handles the full GitHub OAuth login flow with database operations
func HandleGitHubLogin(db *gorm.DB, cfg *GitHubOAuthConfig, code string) (*LoginResponse, error) {
	ctx := context.Background()

	// Exchange code for token
	token, err := cfg.ExchangeCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	// Get user info from GitHub
	githubUser, err := cfg.GetUserInfo(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	// Convert GitHub ID to string
	githubIDStr := fmt.Sprintf("%d", githubUser.ID)

	// Use login name if name is empty
	displayName := githubUser.Name
	if displayName == "" {
		displayName = githubUser.Login
	}

	// Check if user exists by GitHub ID
	var user model.User
	err = db.Where("github_id = ?", githubIDStr).First(&user).Error

	if err == nil {
		// User exists with this GitHub ID - update info and login
		user.AvatarURL = githubUser.AvatarURL
		user.Name = displayName
		user.UpdatedAt = time.Now()
		if err := db.Save(&user).Error; err != nil {
			return nil, fmt.Errorf("failed to update user: %w", err)
		}
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		// Check if user exists with the same email (account linking opportunity)
		if githubUser.Email != "" {
			err = db.Where("email = ?", strings.ToLower(githubUser.Email)).First(&user).Error
			if err == nil {
				// User exists with same email - link GitHub account
				user.GithubID = githubIDStr
				user.AvatarURL = githubUser.AvatarURL
				if user.Name == "" {
					user.Name = displayName
				}
				user.UpdatedAt = time.Now()
				if err := db.Save(&user).Error; err != nil {
					return nil, fmt.Errorf("failed to link GitHub account: %w", err)
				}
			} else {
				// Create new user
				user = model.User{
					ID:        uuid.New().String(),
					Name:      displayName,
					Email:     strings.ToLower(githubUser.Email),
					GithubID:  githubIDStr,
					AvatarURL: githubUser.AvatarURL,
					Status:    "active",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				if err := db.Create(&user).Error; err != nil {
					return nil, fmt.Errorf("failed to create user: %w", err)
				}
			}
		} else {
			// No email from GitHub - use GitHub login as email (not ideal but works)
			user = model.User{
				ID:        uuid.New().String(),
				Name:      displayName,
				Email:     strings.ToLower(githubUser.Login) + "@github.com",
				GithubID:  githubIDStr,
				AvatarURL: githubUser.AvatarURL,
				Status:    "active",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			if err := db.Create(&user).Error; err != nil {
				return nil, fmt.Errorf("failed to create user: %w", err)
			}
		}
	} else {
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Generate JWT token for our app
	appToken, err := GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &LoginResponse{
		Token:     appToken,
		UserID:    user.ID,
		Name:      user.Name,
		Email:     user.Email,
		AvatarURL: user.AvatarURL,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}, nil
}

// GetRedirectURL extracts the frontend URL from the redirect URL for post-login redirect
func GetFrontendURL(redirectURL string) string {
	u, err := url.Parse(redirectURL)
	if err != nil {
		return "/"
	}
	// Extract the base URL (scheme + host) from redirect URL
	// The redirect URL is like http://localhost:8080/api/auth/github/callback
	// We want to redirect back to http://localhost:8080 or the frontend port
	frontendURL := u.Scheme + "://" + u.Host
	return frontendURL
}
