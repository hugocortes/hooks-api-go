// Package middleware provides HTTP middleware
package middleware

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
	"time"

	oidc "github.com/coreos/go-oidc"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

// Middleware provides http middleware
type Middleware struct {
	gin *gin.Engine
}

type oauthConfig struct {
	oauth2   oauth2.Config
	verifier oidc.IDTokenVerifier
}

// New provides middleware funcs
func New(gin *gin.Engine) *Middleware {
	return &Middleware{gin: gin}
}

// NotFound provides 404 route handling
func (h *Middleware) NotFound(c *gin.Context) {
	c.JSON(404, gin.H{"status": 404, "error": "Not Found", "message": "Not Found"})
}

// CorsConfig provides cors
func (h *Middleware) CorsConfig() gin.HandlerFunc {
	// config.AddAllowHeaders("Origin", "X-Requested-With", "Content-Type", "Authorization", "Access-Control-Allow-Origin")
	return cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	})
}

func (h *Middleware) oAuthConfig() oauthConfig {
	realm := os.Getenv("IDP_REALM")
	clientID := os.Getenv("IDP_CLIENT_ID")
	clientSecret := os.Getenv("IDP_CLIENT_SECRET")
	if clientID == "" || clientSecret == "" || realm == "" {
		logrus.Fatal("Missing configuration")
	}

	configURL := os.Getenv("IDP_URI") + "/realms/" + realm
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, configURL)
	if err != nil {
		panic(err)
	}

	oidcConfig := &oidc.Config{
		ClientID: clientID,
	}

	oauth2Config := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID},
	}
	oauth2Config.Endpoint.TokenURL = configURL + "/protocol/openid-connect/token"

	return oauthConfig{
		verifier: *provider.Verifier(oidcConfig),
		oauth2:   oauth2Config,
	}
}

// Auth provides the authentication callback and redirect URL required for
// clients to authenticate themselves
func (h *Middleware) Auth() {
	h.gin.GET("/oauth/redirect", func(c *gin.Context) {
		redirectURI := c.Query("redirect_uri")
		state := c.Query("state")
		if redirectURI == "" || state == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Missing required query params"})
			return
		}
		config := h.oAuthConfig()
		config.oauth2.RedirectURL = redirectURI

		c.Redirect(http.StatusMovedPermanently, config.oauth2.AuthCodeURL(state))
		c.Abort()
		return
	})

	h.gin.POST("/oauth/token", func(c *gin.Context) {
		switch grantType := c.PostForm("grant_type"); grantType {
		case "authorization_code":
			redirectURI := c.PostForm("redirect_uri")
			code := c.PostForm("code")
			if redirectURI == "" || code == "" {
				c.JSON(http.StatusBadRequest, gin.H{"message": "Missing required form params"})
				return
			}

			config := h.oAuthConfig()
			config.oauth2.RedirectURL = redirectURI
			oauthToken, err := config.oauth2.Exchange(c, code)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch token"})
				return
			}
			rawIDToken, ok := oauthToken.Extra("id_token").(string)
			if !ok {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "No id_token field in oauth2 token"})
				return
			}
			_, err = config.verifier.Verify(c, rawIDToken)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to verify token" + err.Error()})
				return
			}

			c.JSON(http.StatusOK, oauthToken)
			return
		case "refresh_token":
			refreshToken := c.PostForm("refresh_token")
			if refreshToken == "" {
				c.JSON(http.StatusBadRequest, gin.H{"message": "Missing refresh token"})
				return
			}

			config := h.oAuthConfig()
			form := url.Values{}
			form.Add("refresh_token", refreshToken)
			form.Add("grant_type", "refresh_token")
			form.Add("scope", strings.Join(config.oauth2.Scopes, " "))
			form.Add("client_id", config.oauth2.ClientID)
			form.Add("client_secret", config.oauth2.ClientSecret)

			client := &http.Client{Timeout: 10 * time.Second}
			req, err := http.NewRequest("POST", config.oauth2.Endpoint.TokenURL, strings.NewReader(form.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal error"})
				return
			}
			resp, err := client.Do(req)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal error"})
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode >= 400 {
				c.JSON(resp.StatusCode, gin.H{"message": "Error"})
				return
			}

			// Token provides the refresh token struct to return
			type Token struct {
				AccessToken  string    `json:"access_token"`
				TokenType    string    `json:"token_type,omitempty"`
				RefreshToken string    `json:"refresh_token,omitempty"`
				Expiry       time.Time `json:"expiry,omitempty"`
				ExpiresIn    int       `json:"expires_in"`
			}
			// Response is the Keycloak response.
			type Response struct {
				Response *http.Response
			}
			// response := &Response{Response: resp}

			// var test = new(Token)
			test := reflect.ValueOf(new(Token)).Interface()
			w, ok := test.(io.Writer)
			if ok {
				io.Copy(w, resp.Body)
			} else {
				json.NewDecoder(resp.Body).Decode(test)
			}
			info, ok := test.(*Token)
			if !ok {
				logrus.Warn("not ok")
			}
			logrus.Debug(info.ExpiresIn)

			// var tokenResponse = new(Token)
			// json.NewDecoder(resp.Body).Decode(tokenResponse)
			// tokenResponse.Expiry = time.Now().Add(time.Duration(tokenResponse.ExpiresIn) * time.Second)

			c.JSON(resp.StatusCode, test)
			return
		default:
			c.JSON(http.StatusBadRequest, gin.H{"message": "Unknown grant type"})
			return
		}
	})
}
