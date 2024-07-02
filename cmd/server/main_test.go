package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

// TestHealthHandler tests the healthHandler function.
func TestHealthHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(healthHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("healthHandler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "OK"
	if rr.Body.String() != expected {
		t.Errorf("healthHandler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

// TestTokenHandler tests the tokenHandler function.
func TestTokenHandler(t *testing.T) {
	tests := []struct {
		name           string
		room           string
		identity       string
		expectedStatus int
	}{
		{
			name:           "Valid request",
			room:           "test-room",
			identity:       "test-identity",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Missing room parameter",
			room:           "",
			identity:       "test-identity",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Missing identity parameter",
			room:           "test-room",
			identity:       "",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("LIVEKIT_API_KEY", "test-api-key")
			os.Setenv("LIVEKIT_API_SECRET", "test-api-secret")

			req, err := http.NewRequest("GET", "/token?room="+tt.room+"&identity="+tt.identity, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(tokenHandler)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("tokenHandler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			if tt.expectedStatus == http.StatusOK {
				responseBody := rr.Body.String()
				if responseBody == "" {
					t.Error("tokenHandler returned empty response")
				}

				// Check if the response is a valid JSON
				var response map[string]string
				err = json.Unmarshal([]byte(responseBody), &response)
				if err != nil {
					t.Errorf("tokenHandler returned invalid JSON: %v", err)
				}

				// Check if the response contains a token and an identity
				if _, ok := response["token"]; !ok {
					t.Error("tokenHandler response does not contain a token")
				}
				if _, ok := response["identity"]; !ok {
					t.Error("tokenHandler response does not contain an identity")
				}
			}
		})
	}
}

// TestCORSAndPreflight tests the CORS headers and preflight requests.
func TestCORSAndPreflight(t *testing.T) {
	tests := []struct {
		name               string
		method             string
		url                string
		expectedStatus     int
		expectedCorsHeader string
	}{
		{
			name:               "Preflight request",
			method:             http.MethodOptions,
			url:                "/token",
			expectedStatus:     http.StatusOK,
			expectedCorsHeader: "*",
		},
		{
			name:               "CORS headers in GET request",
			method:             http.MethodGet,
			url:                "/token?room=test-room&identity=test-identity",
			expectedStatus:     http.StatusOK,
			expectedCorsHeader: "*",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, tt.url, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			var handler http.HandlerFunc
			if tt.url == "/" {
				handler = http.HandlerFunc(healthHandler)
			} else {
				handler = http.HandlerFunc(tokenHandler)
			}

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("%s returned wrong status code: got %v want %v", tt.name, status, tt.expectedStatus)
			}

			if corsHeader := rr.Header().Get("Access-Control-Allow-Origin"); corsHeader != tt.expectedCorsHeader {
				t.Errorf("%s returned wrong CORS header: got %v want %v", tt.name, corsHeader, tt.expectedCorsHeader)
			}
		})
	}
}

// TestMainFunction tests the main function.
func TestMainFunction(t *testing.T) {
	// Set up the required environment variables
	os.Setenv("LIVEKIT_API_KEY", "test-api-key")
	os.Setenv("LIVEKIT_API_SECRET", "test-api-secret")
	os.Setenv("PORT", "8081") // Use a different port to avoid conflicts
	os.Setenv("LOG_LEVEL", "debug")

	// Start the server in a separate goroutine
	go func() {
		main()
	}()

	// Allow some time for the server to start
	time.Sleep(1 * time.Second)

	// Define the test cases
	tests := []struct {
		name           string
		url            string
		expectedStatus int
	}{
		{
			name:           "Health check",
			url:            "http://localhost:8081/",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Valid token request",
			url:            "http://localhost:8081/token?room=test-room&identity=test-identity",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Missing room parameter",
			url:            "http://localhost:8081/token?room=&identity=test-identity",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Missing identity parameter",
			url:            "http://localhost:8081/token?room=test-room&identity=",
			expectedStatus: http.StatusBadRequest,
		},
	}

	// Execute the test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := http.Get(tt.url)
			if err != nil {
				t.Fatalf("Failed to send request: %v", err)
			}
			defer resp.Body.Close()

			if status := resp.StatusCode; status != tt.expectedStatus {
				t.Errorf("Expected status code %v, got %v", tt.expectedStatus, status)
			}
		})
	}

	// Clean up
	os.Unsetenv("LIVEKIT_API_KEY")
	os.Unsetenv("LIVEKIT_API_SECRET")
	os.Unsetenv("PORT")
	os.Unsetenv("LOG_LEVEL")
}
