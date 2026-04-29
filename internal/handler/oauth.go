package handler

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mfb27/luban/internal/auth"
	"github.com/mfb27/luban/internal/response"
	"go.uber.org/zap"
)

const (
	// OAuth state key prefix for Redis storage
	oauthStateKeyPrefix = "oauth:state:"
	// OAuth state expiration time (5 minutes)
	oauthStateExpiration = 5 * time.Minute
)

// These handlers are methods on App struct, not separate type

// GitHubLogin initiates GitHub OAuth flow by returning the authorization URL
func (a *App) githubLoginURL(c *gin.Context) {
	// Check if GitHub OAuth is configured
	if a.cfg.GitHub.ClientID == "" {
		response.NewResponseHelper(c).Error(response.CodeBadRequest, "GitHub OAuth is not configured")
		return
	}

	// Generate state for CSRF protection
	state := generateOAuthState()

	// Store state in Redis with expiration
	ctx := context.Background()
	stateKey := oauthStateKeyPrefix + state
	if err := a.redis.Set(ctx, stateKey, "1", oauthStateExpiration).Err(); err != nil {
		response.NewResponseHelper(c).Error(response.CodeInternal, "Failed to generate OAuth state")
		return
	}

	// Create OAuth config with proxy support
	githubCfg := auth.NewGitHubOAuthConfigWithProxy(a.cfg.GitHub.ClientID, a.cfg.GitHub.ClientSecret, a.cfg.GitHub.RedirectURL, a.cfg.GitHub.ProxyURL)
	authURL := githubCfg.GetAuthURL(state)

	response.NewResponseHelper(c).Success(gin.H{
		"auth_url": authURL,
		"state":    state,
	})
}

// GitHubCallback handles GitHub OAuth callback with code
func (a *App) githubCallback(c *gin.Context) {
	// Get code and state from query parameters
	code := c.Query("code")
	state := c.Query("state")

	if code == "" || state == "" {
		response.NewResponseHelper(c).Error(response.CodeBadRequest, "Missing code or state parameter")
		return
	}

	// Verify state in Redis
	ctx := context.Background()
	stateKey := oauthStateKeyPrefix + state
	val, err := a.redis.Get(ctx, stateKey).Result()
	if err != nil || val == "" {
		response.NewResponseHelper(c).Error(response.CodeAuthFailed, "Invalid or expired OAuth state")
		return
	}

	// Delete state after verification (one-time use)
	a.redis.Del(ctx, stateKey)

	// Check if GitHub OAuth is configured
	if a.cfg.GitHub.ClientID == "" {
		response.NewResponseHelper(c).Error(response.CodeBadRequest, "GitHub OAuth is not configured")
		return
	}

	// Create OAuth config with proxy support
	githubCfg := auth.NewGitHubOAuthConfigWithProxy(a.cfg.GitHub.ClientID, a.cfg.GitHub.ClientSecret, a.cfg.GitHub.RedirectURL, a.cfg.GitHub.ProxyURL)

	// Handle GitHub login
	loginResp, err := auth.HandleGitHubLogin(a.db, githubCfg, code)
	if err != nil {
		a.log.Error("GitHub login failed", zap.Error(err))
		response.NewResponseHelper(c).Error(response.CodeAuthFailed, "GitHub authentication failed")
		return
	}

	// For callback flow, redirect to frontend OAuth callback page
	// The oauth-callback.html will send postMessage to close the popup
	frontendURL := "http://localhost:5501/oauth-callback.html"
	redirectURL := fmt.Sprintf("%s?token=%s&user_id=%s&name=%s", frontendURL, url.QueryEscape(loginResp.Token), url.QueryEscape(loginResp.UserID), url.QueryEscape(loginResp.Name))

	c.Redirect(302, redirectURL)
}

// GitHubLoginDirect handles direct GitHub login with code from frontend (SPA flow)
func (a *App) githubLoginDirect(c *gin.Context) {
	// Check if GitHub OAuth is configured
	if a.cfg.GitHub.ClientID == "" {
		response.NewResponseHelper(c).Error(response.CodeBadRequest, "GitHub OAuth is not configured")
		return
	}

	var input auth.GitHubLoginRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		response.NewResponseHelper(c).Error(response.CodeInvalidParam, err.Error())
		return
	}

	// Create OAuth config with proxy support
	githubCfg := auth.NewGitHubOAuthConfigWithProxy(a.cfg.GitHub.ClientID, a.cfg.GitHub.ClientSecret, a.cfg.GitHub.RedirectURL, a.cfg.GitHub.ProxyURL)

	// Handle GitHub login
	loginResp, err := auth.HandleGitHubLogin(a.db, githubCfg, input.Code)
	if err != nil {
		a.log.Error("GitHub login failed", zap.Error(err))
		response.NewResponseHelper(c).Error(response.CodeAuthFailed, fmt.Sprintf("GitHub authentication failed: %v", err))
		return
	}

	response.NewResponseHelper(c).SuccessWithMessage("GitHub login successful", loginResp)
}

// generateOAuthState generates a random state string
func generateOAuthState() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// LinkGitHubAccount links GitHub account to existing user
func (a *App) linkGitHubAccount(c *gin.Context) {
	// This endpoint requires authentication
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		response.NewResponseHelper(c).Error(response.CodeNoPermission, "User not authenticated")
		return
	}

	// Check if GitHub OAuth is configured
	if a.cfg.GitHub.ClientID == "" {
		response.NewResponseHelper(c).Error(response.CodeBadRequest, "GitHub OAuth is not configured")
		return
	}

	// Generate state for linking
	state := generateOAuthState()
	stateKey := oauthStateKeyPrefix + "link:" + userID.(string)

	ctx := context.Background()
	if err := a.redis.Set(ctx, stateKey, state, oauthStateExpiration).Err(); err != nil {
		response.NewResponseHelper(c).Error(response.CodeInternal, "Failed to generate OAuth state")
		return
	}

	githubCfg := auth.NewGitHubOAuthConfigWithProxy(a.cfg.GitHub.ClientID, a.cfg.GitHub.ClientSecret, a.cfg.GitHub.RedirectURL, a.cfg.GitHub.ProxyURL)
	authURL := githubCfg.GetAuthURL(state)

	response.NewResponseHelper(c).Success(gin.H{
		"auth_url": authURL,
		"state":    state,
	})
}

// GitHubLinkCallback handles GitHub OAuth callback for account linking
func (a *App) githubLinkCallback(c *gin.Context) {
	// This would be similar to githubCallback but for linking accounts
	// For now, we'll use the same callback mechanism
	// The frontend can detect if user is logged in and handle linking vs login
	a.githubCallback(c)
}
