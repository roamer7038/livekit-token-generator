package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/roamer7038/livekit-token-generator/pkg/token"
)

// handler is the HTTP handler for generating join tokens.
func handler(w http.ResponseWriter, r *http.Request) {
	// Set the headers for CORS, content type, and allowed methods
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Content-Type", "application/json")

	// Handle preflight requests
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	room := r.URL.Query().Get("room")
	identity := r.URL.Query().Get("identity")

	if room == "" || identity == "" {
		http.Error(w, "Room and identity parameters are required", http.StatusBadRequest)
		log.Printf("Room and identity parameters are required")
		return
	}

	// Get the VideoGrant from environment variables
	grant := token.GetVideoGrantFromEnv(room)

	// Retrieve the API key and secret from the environment variables
	apiKey := os.Getenv("LIVEKIT_API_KEY")
	apiSecret := os.Getenv("LIVEKIT_API_SECRET")

	params := &token.TokenParams{
		ApiKey:    apiKey,
		ApiSecret: apiSecret,
		Room:      room,
		Identity:  identity,
		Grant:     grant,
	}

	token, err := token.GetJoinToken(params)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		log.Printf("Failed to generate token: %v", err)
		return
	}

	// Create a map to hold the token and identity
	response := map[string]string{
		"token":    token,
		"identity": identity,
	}

	// Convert the map to JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to generate JSON response", http.StatusInternalServerError)
		log.Printf("Failed to generate JSON response: %v", err)
		return
	}

	// Write the JSON response
	w.Write(jsonResponse)

	log.Printf("Token generated for room %s and identity %s", room, identity)
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

// main is the entry point of the application.
func main() {
	// Check if the required environment variables are set
	if os.Getenv("LIVEKIT_API_KEY") == "" || os.Getenv("LIVEKIT_API_SECRET") == "" {
		log.Fatal("LIVEKIT_API_KEY and LIVEKIT_API_SECRET environment variables are required")
	}

	// Retrieve the port from the environment variables
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Retrieve the TLS configuration from the environment variables
	tls := getEnvAsBool("HTTPS", false)

	// Register the handler function
	http.HandleFunc("/", handler)

	ip := "localhost"

	if tls {
		// Check if the required environment variables for TLS are set
		certFile := os.Getenv("SSL_CRT_FILE")
		keyFile := os.Getenv("SSL_KEY_FILE")
		if certFile == "" || keyFile == "" {
			log.Fatal("SSL_CRT_FILE and SSL_KEY_FILE environment variables are required for HTTPS")
		}

		// Start the server with TLS
		log.Printf("Access the server at https://%s:%s", ip, port)
		log.Printf("For example, https://%s:%s?room=room1&identity=user1", ip, port)
		if err := http.ListenAndServeTLS(":"+port, certFile, keyFile, nil); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	} else {
		// Start the server without TLS
		log.Printf("Access the server at http://%s:%s", ip, port)
		log.Printf("For example, http://%s:%s?room=room1&identity=user1", ip, port)
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}
}
