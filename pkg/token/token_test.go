package token

import (
	"os"
	"testing"

	"github.com/livekit/protocol/auth"
)

// TestGetJoinToken tests the GetJoinToken function.
// It checks if the function can generate a non-empty token without errors.
func TestGetJoinToken(t *testing.T) {
	params := &TokenParams{
		ApiKey:    "test-api-key",
		ApiSecret: "test-api-secret",
		Room:      "test-room",
		Identity:  "test-identity",
		Grant:     &auth.VideoGrant{Room: "test-room", RoomJoin: true},
	}

	token, err := GetJoinToken(params)
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

// TestGetVideoGrantFromEnv tests the GetVideoGrantFromEnv function.
// It checks if the function can correctly generate a VideoGrant from environment variables.
func TestGetVideoGrantFromEnv(t *testing.T) {
	os.Setenv("ROOM_CREATE", "true")
	os.Setenv("ROOM_LIST", "false")
	os.Setenv("CAN_PUBLISH", "true")
	os.Setenv("CAN_SUBSCRIBE", "false")

	grant := GetVideoGrantFromEnv("test-room")

	if !grant.RoomCreate || grant.RoomList || !*grant.CanPublish || *grant.CanSubscribe {
		t.Error("Expected values do not match environment variables")
	}

	os.Unsetenv("ROOM_CREATE")
	os.Unsetenv("ROOM_LIST")
	os.Unsetenv("CAN_PUBLISH")
	os.Unsetenv("CAN_SUBSCRIBE")
}
