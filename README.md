# livekit-token-generator

[LiveKit](https://livekit.io/) token generator for Go.

## Requirements

- API Key from LiveKit
- API Secret from LiveKit

## Usage

`cmd/server/main.go` is http server that generates token for LiveKit.

```bash
 env $(cat .env | xargs) go run main.go
```

## License

This project is released under the MIT License. For more details, please refer to the [LICENSE](LICENSE) file.
