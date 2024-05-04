package token

import (
	"os"
	"strconv"
	"time"

	"github.com/livekit/protocol/auth"
)

// TokenParams holds the parameters required to generate a join token.
type TokenParams struct {
	ApiKey    string
	ApiSecret string
	Room      string
	Identity  string
	Grant     *auth.VideoGrant
}

// GetJoinToken generates a join token using the provided parameters.
func GetJoinToken(params *TokenParams) (string, error) {
	at := auth.NewAccessToken(params.ApiKey, params.ApiSecret)
	at.AddGrant(params.Grant).
		SetIdentity(params.Identity).
		SetValidFor(time.Hour)

	return at.ToJWT()
}

// GetVideoGrantFromEnv generates a VideoGrant from environment variables.
func GetVideoGrantFromEnv(room string) *auth.VideoGrant {
	return &auth.VideoGrant{
		RoomCreate:           getEnvAsBool("ROOM_CREATE", false),
		RoomList:             getEnvAsBool("ROOM_LIST", false),
		RoomRecord:           getEnvAsBool("ROOM_RECORD", false),
		RoomAdmin:            getEnvAsBool("ROOM_ADMIN", false),
		RoomJoin:             true,
		Room:                 room,
		CanPublish:           getEnvAsBoolPtr("CAN_PUBLISH", true),
		CanSubscribe:         getEnvAsBoolPtr("CAN_SUBSCRIBE", true),
		CanPublishData:       getEnvAsBoolPtr("CAN_PUBLISH_DATA", true),
		CanUpdateOwnMetadata: getEnvAsBoolPtr("CAN_UPDATE_OWN_METADATA", false),
		IngressAdmin:         getEnvAsBool("INGRESS_ADMIN", false),
		Hidden:               getEnvAsBool("HIDDEN", false),
		Recorder:             getEnvAsBool("RECORDER", false),
		Agent:                getEnvAsBool("AGENT", false),
	}
}

// getEnvAsBool retrieves an environment variable and parses it as a boolean.
// If the environment variable is not set or cannot be parsed, it returns the provided default value.
func getEnvAsBool(name string, defaultVal bool) bool {
	valStr, found := os.LookupEnv(name)
	if !found {
		return defaultVal
	}
	val, err := strconv.ParseBool(valStr)
	if err != nil {
		return defaultVal
	}
	return val
}

// getEnvAsBoolPtr is similar to getEnvAsBool, but returns a pointer to the boolean value.
func getEnvAsBoolPtr(name string, defaultVal bool) *bool {
	valStr, found := os.LookupEnv(name)
	if !found {
		return &defaultVal
	}
	val, err := strconv.ParseBool(valStr)
	if err == nil {
		return &val
	}
	return &defaultVal
}
