# LiveKit Token Generator

This repository contains a Go application that generates tokens for [LiveKit](https://livekit.io/), a real-time video and audio conferencing platform.

## Prerequisites

Before running this application, ensure you have the following:

- LiveKit API Key
- LiveKit API Secret

Both of these can be obtained from your LiveKit account, or you can specify them when self-hosting.

## Quick Start

### Setting Up Environment Variables

1. **Create an Environment File:**

   Start by duplicating the example environment file provided. Rename `.env.example` to `.env` and update the values with your LiveKit credentials:

   ```
   cp .env.example .env
   # Open .env and replace the placeholders with your actual API Key and Secret
   ```

   Your `.env` file should include at least the following variables:

   ```
   LIVEKIT_API_KEY=your_livekit_api_key
   LIVEKIT_API_SECRET=your_livekit_api_secret
   ```

### Running with Docker

To run the server using Docker:

```
docker run -p 8080:8080 --env-file .env ghcr.io/roamer7038/livekit-token-generator
```

This command pulls the image from the GitHub Container Registry and runs it on port 8080, utilizing the environment variables defined in the `.env` file.

### Building and Running the Binary Locally

To build and run the binary locally, navigate to the `cmd/server/` directory and execute the following commands:

```
cd cmd/server/
go build -o livekit-token-generator main.go
env $(cat ../../.env | xargs) ./livekit-token-generator
```

This compiles the main.go file into an executable and runs it, pulling environment variables from your `.env` file.

### How to Use the Application

This application provides a simple HTTP interface to generate LiveKit tokens. To generate a token, you need to specify the room name and the identity (user ID) as URL parameters. Here's how you can request a token:

1. **Request a Token:**

   Use `curl` or any other tool to make an HTTP request. Provide the `room` and `identity` as query parameters. Here are examples of how to do this:

   **Without HTTPS:**

   ```bash
   curl "http://localhost:8080?room=your_room_name&identity=your_identity"
   ```

   **With HTTPS (if configured):**

   ```bash
   curl "https://localhost:8080?room=your_room_name&identity=your_identity"
   ```

   Replace `localhost` with the appropriate IP address or hostname if you are accessing a server deployed elsewhere.

2. **Response:**

   The server will respond with a JSON object containing the `token` and `identity`. Here's an example of the response:

   ```json
   {
     "token": "generated_token_here",
     "identity": "your_identity"
   }
   ```

   Use this token for authenticating with LiveKit to join the specified room.

### Example of Generating a Token

Here is a practical example using `curl`:

```bash
curl "http://localhost:8080?room=DemoRoom&identity=JohnDoe"
```

This command will generate a token for John Doe to join the DemoRoom. You will receive a response in JSON format with the token you need to join the room.

### Error Handling

If required parameters are missing, the server will return an HTTP 400 error. Make sure to include both `room` and `identity` parameters in your request.

## Environment Variables

Server configuration and token settings are controlled by environment variables:

| Variable                | Description                     | Default Value |
| ----------------------- | ------------------------------- | ------------- |
| LIVEKIT_API_KEY         | API Key from LiveKit            | N/A           |
| LIVEKIT_API_SECRET      | API Secret from LiveKit         | N/A           |
| PORT                    | Port number for the HTTP server | 8080          |
| HTTPS                   | Enable HTTPS                    | false         |
| SSL_CRT_FILE            | SSL certificate file path       | N/A           |
| SSL_KEY_FILE            | SSL key file path               | N/A           |
| ROOM_CREATE             | Enable room creation            | false         |
| ROOM_LIST               | Enable room listing             | false         |
| ROOM_RECORD             | Enable room recording           | false         |
| ROOM_ADMIN              | Enable room admin               | false         |
| CAN_PUBLISH             | Enable publishing               | true          |
| CAN_SUBSCRIBE           | Enable subscribing              | true          |
| CAN_PUBLISH_DATA        | Enable data publishing          | true          |
| CAN_UPDATE_OWN_METADATA | Enable updating own metadata    | false         |
| INGRESS_ADMIN           | Enable ingress admin            | false         |
| HIDDEN                  | Enable hidden mode              | false         |
| RECORDER                | Enable recorder function        | false         |
| AGENT                   | Enable agent mode               | false         |

More information on these settings can be found in the [LiveKit documentation](https://docs.livekit.io/realtime/concepts/authentication/).

## License

This project is open-sourced under the MIT License. For more details, see the [LICENSE](LICENSE) file.
