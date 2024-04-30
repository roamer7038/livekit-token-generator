package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/livekit/protocol/auth"
)

// TestGetJoinToken tests the getJoinToken function.
// It checks if the function can generate a non-empty token without errors.
func TestGetJoinToken(t *testing.T) {
	params := &TokenParams{
		ApiKey:    "test-api-key",
		ApiSecret: "test-api-secret",
		Room:      "test-room",
		Identity:  "test-identity",
		Grant:     &auth.VideoGrant{Room: "test-room", RoomJoin: true},
	}

	token, err := getJoinToken(params)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	if token == "" {
		t.Error("Expected token to be non-empty")
	}
}

// TestGetEnvAsBool tests the getEnvAsBool function.
// It checks if the function can correctly parse environment variables as boolean values.
func TestGetEnvAsBool(t *testing.T) {
	os.Setenv("TEST_ENV", "true")
	if !getEnvAsBool("TEST_ENV", false) {
		t.Error("Expected true for environment variable set to 'true'")
	}

	os.Setenv("TEST_ENV", "false")
	if getEnvAsBool("TEST_ENV", true) {
		t.Error("Expected false for environment variable set to 'false'")
	}

	os.Unsetenv("TEST_ENV")
	if getEnvAsBool("TEST_ENV", true) != true {
		t.Error("Expected default value for unset environment variable")
	}
}

// TestGetEnvAsBoolPtr tests the getEnvAsBoolPtr function.
// It checks if the function can correctly parse environment variables as boolean pointers.
func TestGetEnvAsBoolPtr(t *testing.T) {
	os.Setenv("TEST_ENV", "true")
	val := getEnvAsBoolPtr("TEST_ENV", false)
	if val == nil || !*val {
		t.Error("Expected true for environment variable set to 'true'")
	}

	os.Setenv("TEST_ENV", "false")
	val = getEnvAsBoolPtr("TEST_ENV", true)
	if val == nil || *val {
		t.Error("Expected false for environment variable set to 'false'")
	}

	os.Unsetenv("TEST_ENV")
	val = getEnvAsBoolPtr("TEST_ENV", true)
	if val == nil || !*val {
		t.Error("Expected default value for unset environment variable")
	}
}

// TestGetVideoGrantFromEnv tests the getVideoGrantFromEnv function.
// It checks if the function can correctly generate a VideoGrant from environment variables.
func TestGetVideoGrantFromEnv(t *testing.T) {
	os.Setenv("ROOM_CREATE", "true")
	os.Setenv("ROOM_LIST", "false")
	os.Setenv("CAN_PUBLISH", "true")
	os.Setenv("CAN_SUBSCRIBE", "false")

	grant := getVideoGrantFromEnv("test-room")

	if !grant.RoomCreate || grant.RoomList || !*grant.CanPublish || *grant.CanSubscribe {
		t.Error("Expected values do not match environment variables")
	}

	os.Unsetenv("ROOM_CREATE")
	os.Unsetenv("ROOM_LIST")
	os.Unsetenv("CAN_PUBLISH")
	os.Unsetenv("CAN_SUBSCRIBE")
}

// TestHandler tests the handler function.
// It checks if the function can correctly handle an HTTP request and return a non-empty response with a 200 status code.
func TestHandler(t *testing.T) {
	os.Setenv("LIVEKIT_API_KEY", "test-api-key")
	os.Setenv("LIVEKIT_API_SECRET", "test-api-secret")

	req, err := http.NewRequest("GET", "/?room=test-room&identity=test-identity", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

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