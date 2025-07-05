# Go LLM-Gateway

## Overview
go-llm-gateway is a gateway server that proxies chat completion requests to multiple AI providers (such as OpenAI and Anthropic) through a unified API. It supports health checking of provider endpoints and secure routing based on virtual API keys.

---

## Setup Instructions

### Prerequisites
- Go 1.18 or newer

### Installation
1. Clone the repository:
   ```bash
   git clone <repo-url>
   cd go-llm-gateway
   ```
2. Install dependencies:
   ```bash
   go mod tidy
   ```
3. Build and run the server:
   ```bash
   go run ./cmd/server/server.go
   ```

4. The server will start and listen on port 8080.
5. You may use the client to test the server:
   ```bash
   go run ./cmd/client/client.go
   ```

   * See the -h flag for more information.
---

## Configuration

### API Keys
- `keys.json` contains your virtual API keys and their mapping to providers. Example:
  ```json
  {
    "virtual_keys": {
    "vk_user1_openai": {
      "provider": "openai",
      "api_key": "sk-real-openai-key-123"
    },
    "vk_user2_anthropic": {
      "provider": "anthropic",
      "api_key": "sk-ant-real-anthropic-key-456"
    },
    "vk_admin_openai": {
      "provider": "openai",
      "api_key": "sk-another-openai-key-789"
    }
  }
  ```

### Provider Endpoints
- `internal/provider/aiservices.json` defines available AI service endpoints:
  ```json
  {
    "openai": "https://api.openai.com/v1/chat/completions",
    "anthropic": "https://api.anthropic.com/v1/messages"
  }
  ```
---

## Usage

### Health API
- **Endpoint:** `GET /health`
- **Description:** Returns the health status of all configured AI providers.
- **Example Response:**
  ```json
  {
    "status": "ok",
    "providers": [
      {"name": "openai", "url": "https://api.openai.com/v1/chat/completions", "online": true},
      {"name": "anthropic", "url": "https://api.anthropic.com/v1/messages", "online": false}
    ]
  }
  ```

### Chat Completions API
- **Endpoint:** `POST /completions/chat`
- **Headers:**
  - `Authorization: Bearer <virtual-key>`
- **Body:** (forwarded as-is to the provider)
  ```json
  {
    "model": "gpt-3.5-turbo",
    "messages": [
      {"role": "user", "content": "Hello, how are you?"}
    ]
  }
  ```
- **How Provider Routing Works:**
  - The server reads the `virtual-key` from the `Authorization` header.
  - It looks up the provider and API key in `keys.json`.
  - It finds the provider endpoint in `internal/provider/aiservices.json`.
  - The request is proxied to the provider with the correct API key.

#### Example Request
```bash
curl -X POST http://localhost:8080/completions/chat \
  -H 'Authorization: Bearer testkey' \
  -H 'Content-Type: application/json' \
  -d '{
    "model": "gpt-3.5-turbo",
    "messages": [
      {"role": "user", "content": "Hello, how are you?"}
    ]
  }'
```

#### Example Response
```json
{
  "choices": [
    {
      "message": {
        "role": "assistant",
        "content": "Hello! I'm doing well, thank you for asking."
      }
    }
  ]
}
```

#### Error Cases
- Invalid or missing Authorization header:
  - Response: `400 Bad Request`
- Unknown provider in `keys.json`:
  - Response: `400 Bad Request`
- Unauthorized key:
  - Response: `401 Unauthorized`

---

## Running Tests
Run all tests with:
```bash
go test ./...
```

---

## Project Structure
- `cmd/server/server.go` — Main server entrypoint
- `internal/handler/` — HTTP handlers for endpoints
- `internal/provider/` — Provider configuration and endpoints
- `internal/config/` — Configuration loading
