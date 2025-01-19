package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"rest-api/internal/models"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

type OAuthService struct {
	config  *oauth2.Config
	userSvc *UserService
	authSvc *AuthService
}

func NewOAuthService(clientID, clientSecret, redirectURL string, userSvc *UserService, authSvc *AuthService) *OAuthService {
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return &OAuthService{
		config:  config,
		userSvc: userSvc,
		authSvc: authSvc,
	}
}

func (s *OAuthService) GetAuthURL(state string) string {
	return s.config.AuthCodeURL(state)
}

func (s *OAuthService) HandleCallback(code string) (*models.Session, error) {
	token, err := s.config.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %w", err)
	}

	userInfo, err := s.getUserInfo(token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	// Check if user exists
	user, err := s.userSvc.FindByEmail(userInfo.Email)
	if err != nil {
		// Create new user if not exists
		user, err = s.userSvc.CreateGoogleUser(userInfo.Email, userInfo.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}
	}

	// Create session
	session, err := s.authSvc.CreateSession(user, "Google OAuth")
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return session, nil
}

func (s *OAuthService) getUserInfo(accessToken string) (*GoogleUserInfo, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + accessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %w", err)
	}

	var userInfo GoogleUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, fmt.Errorf("failed parsing user info: %w", err)
	}

	return &userInfo, nil
}
