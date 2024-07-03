package main

import (
	"encoding/json"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/roamer7038/livekit-token-generator/pkg/token"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// healthHandler handles health check requests.
func healthHandler(w http.ResponseWriter, r *http.Request) {
	// Set the headers for CORS, content type, and allowed methods
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// Handle preflight requests
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
	log.Debug().Str("method", r.Method).Str("path", r.URL.Path).Str("remote_addr", r.RemoteAddr).Msg("Health check")
}

// tokenHandler handles token generation requests.
func tokenHandler(w http.ResponseWriter, r *http.Request) {
	// Set the headers for CORS, content type, and allowed methods
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// Handle preflight requests
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	room := r.URL.Query().Get("room")
	identity := r.URL.Query().Get("identity")

	if room == "" || identity == "" {
		http.Error(w, "Room and identity parameters are required", http.StatusBadRequest)
		log.Warn().Msg("Room and identity parameters are required")
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

	joinToken, err := token.GetJoinToken(params)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		log.Error().Err(err).Msg("Failed to generate token")
		return
	}

	// Create a map to hold the token and identity
	response := map[string]string{
		"token":    joinToken,
		"identity": identity,
		"room":     room,
	}

	// Convert the map to JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to generate JSON response", http.StatusInternalServerError)
		log.Error().Err(err).Msg("Failed to generate JSON response")
		return
	}

	// Set content type before writing the response
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)

	log.Info().Str("method", r.Method).Str("path", r.URL.Path).Str("remote_addr", r.RemoteAddr).Str("room", room).Str("identity", identity).Msg("Token generated")
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

// getLocalIP retrieves the local IP address of the server.
func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "unknown"
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String()
			}
		}
	}
	return "unknown"
}

// main is the entry point of the application.
func main() {
	// Set log level based on environment variable
	logLevel := os.Getenv("LOG_LEVEL")
	switch logLevel {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	apiKey := os.Getenv("LIVEKIT_API_KEY")
	apiSecret := os.Getenv("LIVEKIT_API_SECRET")
	if apiKey == "" || apiSecret == "" {
		log.Fatal().Msg("LIVEKIT_API_KEY and LIVEKIT_API_SECRET environment variables are required")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	tls := getEnvAsBool("HTTPS", false)

	http.HandleFunc("/", healthHandler)
	http.HandleFunc("/token", tokenHandler)

	ip := getLocalIP()
	log.Debug().Msgf("Local IP address: %s", ip)

	if tls {
		certFile := os.Getenv("SSL_CRT_FILE")
		keyFile := os.Getenv("SSL_KEY_FILE")
		if certFile == "" || keyFile == "" {
			log.Fatal().Msg("SSL_CRT_FILE and SSL_KEY_FILE environment variables are required for HTTPS")
		}

		log.Info().Msgf("Access the server at https://%s:%s", ip, port)
		log.Info().Msgf("For example, https://%s:%s/token?room=room1&identity=user1", ip, port)
		if err := http.ListenAndServeTLS(":"+port, certFile, keyFile, nil); err != nil {
			log.Fatal().Err(err).Msg("Failed to start server")
		}
	} else {
		log.Info().Msgf("Access the server at http://%s:%s", ip, port)
		log.Info().Msgf("For example, http://%s:%s/token?room=room1&identity=user1", ip, port)
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			log.Fatal().Err(err).Msg("Failed to start server")
		}
	}
}
