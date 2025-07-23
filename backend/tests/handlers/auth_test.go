package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"codewithdell/backend/internal/handlers"
	"codewithdell/backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockDB is a mock database for testing
type MockDB struct {
	mock.Mock
}

func (m *MockDB) Where(query interface{}, args ...interface{}) *gorm.DB {
	mockArgs := m.Called(query, args)
	return mockArgs.Get(0).(*gorm.DB)
}

func (m *MockDB) First(dest interface{}, conds ...interface{}) *gorm.DB {
	mockArgs := m.Called(dest, conds)
	return mockArgs.Get(0).(*gorm.DB)
}

func (m *MockDB) Create(value interface{}) *gorm.DB {
	mockArgs := m.Called(value)
	return mockArgs.Get(0).(*gorm.DB)
}

func (m *MockDB) Error() error {
	mockArgs := m.Called()
	return mockArgs.Error(0)
}

// TestRegister tests user registration
func TestRegister(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    handlers.RegisterRequest
		expectedStatus int
		expectedError  string
		setupMock      func(*MockDB)
	}{
		{
			name: "Valid registration",
			requestBody: handlers.RegisterRequest{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john@example.com",
				Password:  "Password123!",
				Username:  "johndoe",
			},
			expectedStatus: http.StatusCreated,
			setupMock: func(mockDB *MockDB) {
				// Mock user existence check (user doesn't exist)
				mockDB.On("Where", "email = ? OR username = ?", "john@example.com", "johndoe").
					Return(mockDB)
				mockDB.On("First", mock.Anything, mock.Anything).
					Return(mockDB)
				mockDB.On("Error").Return(gorm.ErrRecordNotFound)

				// Mock user creation
				mockDB.On("Create", mock.Anything).
					Return(mockDB)
				mockDB.On("Error").Return(nil)
			},
		},
		{
			name: "User already exists",
			requestBody: handlers.RegisterRequest{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "existing@example.com",
				Password:  "Password123!",
				Username:  "existinguser",
			},
			expectedStatus: http.StatusConflict,
			expectedError:  "User already exists",
			setupMock: func(mockDB *MockDB) {
				// Mock user existence check (user exists)
				mockDB.On("Where", "email = ? OR username = ?", "existing@example.com", "existinguser").
					Return(mockDB)
				mockDB.On("First", mock.Anything, mock.Anything).
					Return(mockDB)
				mockDB.On("Error").Return(nil)
			},
		},
		{
			name: "Invalid email",
			requestBody: handlers.RegisterRequest{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "invalid-email",
				Password:  "Password123!",
				Username:  "johndoe",
			},
			expectedStatus: http.StatusBadRequest,
			setupMock:      func(mockDB *MockDB) {},
		},
		{
			name: "Weak password",
			requestBody: handlers.RegisterRequest{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john@example.com",
				Password:  "weak",
				Username:  "johndoe",
			},
			expectedStatus: http.StatusBadRequest,
			setupMock:      func(mockDB *MockDB) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockDB := &MockDB{}
			tt.setupMock(mockDB)

			router := gin.New()
			router.Use(func(c *gin.Context) {
				c.Set("db", mockDB)
				c.Next()
			})

			router.POST("/register", handlers.Register)

			// Create request
			jsonBody, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError != "" {
				var response map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &response)
				assert.Contains(t, response["error"], tt.expectedError)
			}

			if tt.expectedStatus == http.StatusCreated {
				var response handlers.AuthResponse
				json.Unmarshal(w.Body.Bytes(), &response)
				assert.NotEmpty(t, response.Token)
				assert.NotEmpty(t, response.RefreshToken)
				assert.Equal(t, tt.requestBody.Email, response.User.Email)
				assert.Equal(t, tt.requestBody.Username, response.User.Username)
			}

			mockDB.AssertExpectations(t)
		})
	}
}

// TestLogin tests user login
func TestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
	name           string
	requestBody    handlers.AuthRequest
	expectedStatus int
	expectedError  string
	setupMock      func(*MockDB)
}{
	{
		name: "Valid login",
		requestBody: handlers.AuthRequest{
			Email:    "john@example.com",
			Password: "Password123!",
		},
		expectedStatus: http.StatusOK,
		setupMock: func(mockDB *MockDB) {
			// Mock user lookup
			mockDB.On("Where", "email = ?", "john@example.com").
				Return(mockDB)
			mockDB.On("First", mock.Anything, mock.Anything).
				Return(mockDB)
			mockDB.On("Error").Return(nil)
		},
	},
	{
		name: "Invalid credentials",
		requestBody: handlers.AuthRequest{
			Email:    "nonexistent@example.com",
			Password: "wrongpassword",
		},
		expectedStatus: http.StatusUnauthorized,
		expectedError:  "Invalid credentials",
		setupMock: func(mockDB *MockDB) {
			// Mock user lookup (user not found)
			mockDB.On("Where", "email = ?", "nonexistent@example.com").
				Return(mockDB)
			mockDB.On("First", mock.Anything, mock.Anything).
				Return(mockDB)
			mockDB.On("Error").Return(gorm.ErrRecordNotFound)
		},
	},
	{
		name: "Invalid email format",
		requestBody: handlers.AuthRequest{
			Email:    "invalid-email",
			Password: "Password123!",
		},
		expectedStatus: http.StatusBadRequest,
		setupMock:      func(mockDB *MockDB) {},
	},
}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockDB := &MockDB{}
			tt.setupMock(mockDB)

			router := gin.New()
			router.Use(func(c *gin.Context) {
				c.Set("db", mockDB)
				c.Next()
			})

			router.POST("/login", handlers.Login)

			// Create request
			jsonBody, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError != "" {
				var response map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &response)
				assert.Contains(t, response["error"], tt.expectedError)
			}

			if tt.expectedStatus == http.StatusOK {
				var response handlers.AuthResponse
				json.Unmarshal(w.Body.Bytes(), &response)
				assert.NotEmpty(t, response.Token)
				assert.NotEmpty(t, response.RefreshToken)
			}

			mockDB.AssertExpectations(t)
		})
	}
}

// TestRefreshToken tests token refresh
func TestRefreshToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		refreshToken   string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "Missing refresh token",
			refreshToken:   "",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Refresh token required",
		},
		{
			name:           "Invalid refresh token",
			refreshToken:   "invalid-token",
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Invalid refresh token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.POST("/refresh", handlers.RefreshToken)

			// Create request
			req, _ := http.NewRequest("POST", "/refresh", nil)
			if tt.refreshToken != "" {
				req.Header.Set("Authorization", "Bearer "+tt.refreshToken)
			}

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError != "" {
				var response map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &response)
				assert.Contains(t, response["error"], tt.expectedError)
			}
		})
	}
} 