package e2e_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	baseURL        string
	userToken      string
	refreshToken   string
	adminToken     string
	librarianToken string
	userId         int
	authorId       int
	bookId         int
}

type AuthResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int    `json:"expiresIn"`
	TokenType    string `json:"tokenType"`
	UserId       int    `json:"userId"`
	Username     string `json:"username"`
	Role         string `json:"role"`
}

type ErrorResponse struct {
	Message string   `json:"message"`
	Errors  []string `json:"errors,omitempty"`
}

func (suite *TestSuite) SetupSuite() {
	err := godotenv.Load("../../.env")
	if err != nil {
		suite.T().Logf("Warning: .env file not found, using default or system environment variables: %v", err)
	}

	suite.baseURL = getEnv("API_BASE_URL", "http://localhost:8080/api/v1")
	suite.registerTestUser()
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func (suite *TestSuite) sendRequest(method, endpoint string, payload interface{}, token string) (*http.Response, error) {
	var reqBody io.Reader

	if payload != nil {
		jsonData, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal JSON: %v", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, suite.baseURL+endpoint, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	if reqBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	return client.Do(req)
}

func (suite *TestSuite) registerTestUser() {
	testUser := map[string]interface{}{
		"username":   "admin1",
		"email":      "admin1@example.com",
		"password":   "admin1",
		"first_name": "Admin",
		"last_name":  "User",
		"role":       "admin",
	}

	resp, err := suite.sendRequest("POST", "/auth/register", testUser, "")
	if err != nil {
		suite.T().Fatalf("Failed to register test user: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 && resp.StatusCode != 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		suite.T().Fatalf("Failed to register test user. Status: %d, Response: %s", resp.StatusCode, string(bodyBytes))
	}
}

func (suite *TestSuite) TestUserLogin() {
	t := suite.T()

	loginPayload := map[string]interface{}{
		"username": "admin1",
		"password": "admin1",
	}

	resp, err := suite.sendRequest("POST", "/auth/login", loginPayload, "")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var authResp AuthResponse
	err = json.NewDecoder(resp.Body).Decode(&authResp)
	require.NoError(t, err)

	assert.NotEmpty(t, authResp.Token)
	assert.NotEmpty(t, authResp.RefreshToken)
	assert.NotZero(t, authResp.ExpiresIn)
	assert.Equal(t, "Bearer", authResp.TokenType)
	assert.NotZero(t, authResp.UserId)
	assert.Equal(t, "admin1", authResp.Username)
	assert.Equal(t, "admin", authResp.Role)

	suite.userToken = authResp.Token
	suite.refreshToken = authResp.RefreshToken
	suite.userId = authResp.UserId
}

func (suite *TestSuite) TestRefreshToken() {
	t := suite.T()

	require.NotEmpty(t, suite.refreshToken)

	refreshPayload := map[string]interface{}{
		"refresh_token": suite.refreshToken,
	}

	resp, err := suite.sendRequest("POST", "/auth/refresh", refreshPayload, "")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var authResp AuthResponse
	err = json.NewDecoder(resp.Body).Decode(&authResp)
	require.NoError(t, err)

	assert.NotEmpty(t, authResp.Token)
	assert.NotEmpty(t, authResp.RefreshToken)
	assert.NotZero(t, authResp.ExpiresIn)
	assert.Equal(t, "Bearer", authResp.TokenType)

	suite.userToken = authResp.Token
	suite.refreshToken = authResp.RefreshToken
}

func (suite *TestSuite) TestLoginWithInvalidCredentials() {
	t := suite.T()

	loginPayload := map[string]interface{}{
		"username": "admin1",
		"password": "wrongpassword",
	}

	resp, err := suite.sendRequest("POST", "/auth/login", loginPayload, "")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var errResp ErrorResponse
	err = json.NewDecoder(resp.Body).Decode(&errResp)
	require.NoError(t, err)

	assert.Contains(t, errResp.Message, "invalid credentials")
}

func TestAuthenticationSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
