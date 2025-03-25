# Real-time Translation Service

Real-time translation service using Azure Speech Service and Azure Translator.

日本語版のREADMEは[こちら](./README-ja.md)をご覧ください。

## Generating Go Client Library with AutoRest

We use [AutoRest](https://github.com/Azure/autorest) to generate a client library by loading API specifications. The approach allows generating client libraries in any supported language.

For this project, we load the Translator API specification to generate a Go client library:

```bash
autorest --go --input-file=https://raw.githubusercontent.com/Azure/azure-rest-api-specs/refs/heads/master/specification/cognitiveservices/data-plane/TranslatorText/stable/v3.0/TranslatorText.json --output-folder=./translatortext --namespace=translatortext
```

By adding the generated module to the project, you can make API requests to Azure resources.

- [Translator API Spec v3.0](https://learn.microsoft.com/en-us/azure/ai-services/translator/text-translation/reference/v3/reference)

### Issues

- When receiving translation result responses, the following error might occur. Here's how to fix it:

```
unmarshalling type *[]*translatortext.TranslateResultAllItem: unmarshalling type *translatortext.TranslateResultAllItem: struct field DetectedLanguage: unmarshalling type *translatortext.TranslateResultAllItemDetectedLanguage: struct field Score: json: cannot unmarshal number 1.0 into Go value of type int32
→
// Find the struct like this:
type TranslateResultAllItemDetectedLanguage struct {
    Language string `json:"language,omitempty"`
    Score    int32  `json:"score,omitempty"` // This field type is the issue
}
 
// Change it to:
type TranslateResultAllItemDetectedLanguage struct {
    Language string  `json:"language,omitempty"`
    Score    float64 `json:"score,omitempty"` // Changed from int32 to float64
}
```

## Setup Instructions

- [How to setup backend API](./backend/README.md)

When the server starts successfully, you will see a message like:
```
Speech Recognition and Translation Server is running on port 8080
```

- [How to setup frontend App](./frontend/README.md)

## Environment Variables

The application requires the following environment variables:

### Backend:
- `AZURE_CLIENT_ID` - Service Principal Client ID
- `AZURE_CLIENT_SECRET` - Service Principal Secret
- `AZURE_TENANT_ID` - Entra ID Tenant ID
- `SPEECH_SERVICE_KEY` - Azure Speech Service subscription key
- `SPEECH_SERVICE_REGION` - Azure Speech Service region (e.g., japaneast)
- `PORT` - Port used by the server (default: 8080)

## Stopping the API

To stop the server, press `Ctrl+C` in the terminal. A graceful shutdown will be performed.

## API Usage Examples with curl

You can interact with the translation API using curl commands as follows:

### Text Translation

To translate text from one language to another:

```bash
curl -X POST http://localhost:8080/api/v1/translate \
  -H "Content-Type: application/json" \
  -d '{
    "text": "Hello, how are you?",
    "sourceLanguage": "en",
    "targetLanguage": "ja"
  }'
```

### Streaming Translation

#### 1. Start a streaming session

```bash
curl -X POST http://localhost:8080/api/v1/streaming/start \
  -H "Content-Type: application/json" \
  -d '{
    "sourceLanguage": "en",
    "targetLanguage": "ja",
    "audioFormat": "wav"
  }'
```

Response will include a `sessionId` and `webSocketURL` that you'll need for subsequent requests:

```json
{
  "sessionId": "12345678-1234-1234-1234-123456789abc",
  "webSocketURL": "/api/v1/streaming/ws/12345678-1234-1234-1234-123456789abc",
  "sourceLanguage": "en",
  "targetLanguage": "ja"
}
```

#### 2. Connect to WebSocket for real-time translation

For real-time translation, connect to the WebSocket URL provided in the response to the start endpoint.
After connecting, follow the WebSocket protocol documented in the backend README.

#### 3. Process audio chunks (Alternative to WebSocket)

```bash
curl -X POST http://localhost:8080/api/v1/streaming/process \
  -H "Content-Type: application/json" \
  -d '{
    "sessionId": "12345678-1234-1234-1234-123456789abc",
    "audioChunk": "BASE64_ENCODED_AUDIO_DATA"
  }'
```

#### 4. Close the streaming session

```bash
curl -X POST http://localhost:8080/api/v1/streaming/close \
  -H "Content-Type: application/json" \
  -d '{
    "sessionId": "12345678-1234-1234-1234-123456789abc"
  }'
```

### Health Check

To check if the API server is running:

```bash
curl http://localhost:8080/api/v1/health
```

Expected response:

```json
{
  "status": "ok"
}
```