package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockStorage реализует интерфейс Storage для тестов
type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) Get(key string) *string {
	args := m.Called(key)
	if v, ok := args.Get(0).(string); ok {
		return &v
	}
	return nil
}

func (m *MockStorage) Set(key, value string) {
	m.Called(key, value)
}

// Test for /health endpoint
func TestHealthEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockStorage := new(MockStorage)
	server := NewServer(mockStorage)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/health", nil)
	server.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "OK", response["status"])
}

// Test for GET /scalar/get/:key endpoint
func TestGetScalarEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockStorage := new(MockStorage)
	server := NewServer(mockStorage)

	key := "testKey"
	value := "testValue"
	mockStorage.On("Get", key).Return(value)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/scalar/get/"+key, nil)
	server.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, key, response["key"])
	assert.Equal(t, value, response["value"])
}

// Test for GET /scalar/get/:key with a missing key
func TestGetScalarEndpointNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockStorage := new(MockStorage)
	server := NewServer(mockStorage)

	key := "missingKey"
	mockStorage.On("Get", key).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/scalar/get/"+key, nil)
	server.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "key not found", response["error"])
}

// Test for PUT /scalar/set/:key endpoint
func TestSetScalarEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockStorage := new(MockStorage)
	server := NewServer(mockStorage)

	key := "testKey"
	value := "testValue"
	mockStorage.On("Set", key, value).Return()

	w := httptest.NewRecorder()
	body := map[string]string{"Value": value}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPut, "/scalar/set/"+key, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	server.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockStorage.AssertCalled(t, "Set", key, value)
}

// Test for PUT /scalar/set/:key with invalid JSON body
func TestSetScalarEndpointInvalidBody(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockStorage := new(MockStorage)
	server := NewServer(mockStorage)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/scalar/set/testKey", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	server.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "invalid request body", response["error"])
}
