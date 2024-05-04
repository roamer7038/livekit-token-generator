package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// TestHandler tests the handler function.
func TestHandler(t *testing.T) {
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

			req, err := http.NewRequest("GET", "/?room="+tt.room+"&identity="+tt.identity, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(handler)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			if tt.expectedStatus == http.StatusOK {
				responseBody := rr.Body.String()
				if responseBody == "" {
					t.Error("handler returned empty response")
				}

				// Check if the response is a valid JSON
				var response map[string]string
				err = json.Unmarshal([]byte(responseBody), &response)
				if err != nil {
					t.Errorf("handler returned invalid JSON: %v", err)
				}

				// Check if the response contains a token and an identity
				if _, ok := response["token"]; !ok {
					t.Error("handler response does not contain a token")
				}
				if _, ok := response["identity"]; !ok {
					t.Error("handler response does not contain an identity")
				}
			}
		})
	}
}
