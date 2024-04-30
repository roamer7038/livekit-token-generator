package main

import (
	"fmt"
	"log"
	"net/http"
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

// getJoinToken generates a join token using the provided parameters.
func getJoinToken(params *TokenParams) (string, error) {
	at := auth.NewAccessToken(params.ApiKey, params.ApiSecret)
	at.AddGrant(params.Grant).
		SetIdentity(params.Identity).
		SetValidFor(time.Hour)

	return at.ToJWT()
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

// getVideoGrantFromEnv generates a VideoGrant from environment variables.
func getVideoGrantFromEnv(room string) *auth.VideoGrant {
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

// handler is the HTTP handler for generating join tokens.
func handler(w http.ResponseWriter, r *http.Request) {
	room := r.URL.Query().Get("room")
	identity := r.URL.Query().Get("identity")

	if room == "" || identity == "" {
		http.Error(w, "Room and identity parameters are required", http.StatusBadRequest)
		return
	}

	grant := getVideoGrantFromEnv(room)

	apiKey := os.Getenv("LIVEKIT_API_KEY")
	apiSecret := os.Getenv("LIVEKIT_API_SECRET")

	params := &TokenParams{
		ApiKey:    apiKey,
		ApiSecret: apiSecret,
		Room:      room,
		Identity:  identity,
		Grant:     grant,
	}

	token, err := getJoinToken(params)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		log.Printf("Failed to generate token: %v", err)
		return
	}

	fmt.Fprint(w, token)
}

// main is the entry point of the application.
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	tls := getEnvAsBool("HTTPS", false)
	certFile := os.Getenv("SSL_CRT_FILE")
	keyFile := os.Getenv("SSL_KEY_FILE")

	http.HandleFunc("/", handler)

	log.Printf("Server starting on port %s", port)
	if tls {
		if err := http.ListenAndServeTLS(":"+port, certFile, keyFile, nil); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	} else {
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}
}
