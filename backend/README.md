# Real-time Speech Translation Service Backend

## Overview

This backend service is a RESTful API that provides text translation and real-time speech streaming translation functionality using Azure Translator Service. It is implemented in Go and uses the Gin framework.

## System Architecture

The system consists of the following components:

- **Gin Web Server**: Processes HTTP requests and provides various endpoints
- **Azure Translator Client**: Client for communicating with the Azure Translator Text API
- **Session Management**: In-memory storage for managing streaming translation sessions
- **Audio Processing**: Module for processing Base64 encoded audio data

```
+----------------+        +-------------------+
|                |        |                   |
| Client         +------->+ Gin Web Server    |
|                |        |                   |
+----------------+        +--------+----------+
                                  |
                                  v
                          +----------------+         +------------------+
                          |                |         |                  |
                          | TranslatorClient+-------->+ Azure Translator |
                          |                |         |                  |
                          +----------------+         +------------------+
```

## API Endpoints

### Health Check

```
GET /api/v1/health
```

Endpoint to check the server status.

**Response Example**:
```json
{
  "status": "ok"
}
```

### Text Translation

```
POST /api/v1/translate
```

Translates text to the specified language.

**Request Example**:
```json
{
  "text": "こんにちは",
  "targetLanguage": "en",
  "sourceLanguage": "ja"
}
```

**Response Example**:
```json
{
  "originalText": "こんにちは",
  "translatedText": "Hello",
  "sourceLanguage": "ja",
  "targetLanguage": "en",
  "confidence": 0.98
}
```

### Start Streaming Translation Session

```
POST /api/v1/streaming/start
```

Starts a streaming translation session.

**Request Example**:
```json
{
  "sourceLanguage": "ja",
  "targetLanguage": "en",
  "audioFormat": "wav"
}
```

**Response Example**:
```json
{
  "sessionId": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
}
```

### Process Audio Data

```
POST /api/v1/streaming/process
```

Sends and processes Base64 encoded audio chunks.

**Request Example**:
```json
{
  "sessionId": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "audioChunk": "UklGRjoAAABXQVZFZm10IBIAAAAHAAEAQB8AAEAfAAABAAgAAABMSVNUHAAAAElORk9JU0ZUDQAAAExhdmY1OC4yOS4xMDDA/w=="
}
```

**Response Example**:
```json
[
  {
    "sourceLanguage": "ja",
    "targetLanguage": "en",
    "translatedText": "Hello, how are you?",
    "isFinal": true,
    "segmentId": "f7e8d9c0-b1a2-3456-7890-abcdef123456"
  }
]
```

### Close Streaming Session

```
POST /api/v1/streaming/close
```

Ends a streaming session.

**Request Example**:
```json
{
  "sessionId": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
}
```

**Response Example**:
```json
{
  "status": "Session closed"
}
```

## Audio Data Requirements

- Supported formats: WAV
- Sampling rate: 16kHz recommended
- Bit depth: 16bit
- Channels: Mono
- Base64 encoding: Audio data must be sent Base64 encoded

## Creating a Service Principal

Create a service principal using Azure CLI.

```bash
az ad sp create-for-rbac --name "go-translation-service" --role contributor --scopes /subscriptions/{subscription-id}/resourceGroups/{resource-group}
```

After executing the command, the following information will be displayed:
- appId (AZURE_CLIENT_ID)
- password (AZURE_CLIENT_SECRET)
- tenant (AZURE_TENANT_ID)

## Setting Permissions

- For simplicity, grant `Contributor` at the resource group scope.
- In production environments, it is recommended to follow the principle of least privilege and grant only necessary permissions.

## Setting Environment Variables

- Copy the `.env.example` file and create a `.env` file.

```bash
cp .env.example .env
```

- Set the following environment variables in the `.env` file.

| Environment Variable | Description |
|----------|------|
| AZURE_CLIENT_ID | Service Principal Client ID |
| AZURE_CLIENT_SECRET | Service Principal Secret |
| AZURE_TENANT_ID | Entra ID Tenant ID |
| TRANSLATOR_SUBSCRIPTION_KEY | Azure Translator resource subscription key |
| TRANSLATOR_SUBSCRIPTION_REGION | Azure Translator resource region (e.g., japaneast) |
| PORT | Port used by the server (default: 8080) |

## Local Development

### Requirements

- Go 1.16 or higher
- Azure subscription
- Azure Translator resource

### Local Execution

```bash
go run main.go
```

## Running with Docker

```bash
# Build Docker image
docker build -t go-translation-service .

# Run container
docker run --env-file .env -p 8080:8080 go-translation-service
```

## Error Handling

The service returns the following HTTP status codes:

- 200 OK: Request successful
- 400 Bad Request: Invalid request parameters
- 401 Unauthorized: Authentication failed
- 404 Not Found: Resource not found
- 500 Internal Server Error: Server internal error

## Performance Considerations

- Streaming sessions are managed in memory, so all sessions will be lost when the server restarts
- For large-scale environments, consider using external caching like Redis to store session state
- Consider implementing a timeout mechanism to automatically delete sessions that have been idle for a long time

## Supported Languages

For a list of supported languages, refer to the Azure Translator Service documentation. Currently, more than 100 languages are supported.