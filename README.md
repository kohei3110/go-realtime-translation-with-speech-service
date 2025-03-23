# Real-time Translation Service

Real-time translation service using Azure Speech Service.

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
Server is running on port 8080
```

- [How to setup frontend App](./frontend/README.md)

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
    "audioFormat": "audio/wav"
  }'
```

Response will include a `sessionId` that you'll need for subsequent requests:
```json
{
  "sessionId": "12345678-1234-1234-1234-123456789abc"
}
```

#### 2. Process audio chunks

```bash
curl -X POST http://localhost:8080/api/v1/streaming/process \
  -H "Content-Type: application/json" \
  -d '{
    "sessionId": "12345678-1234-1234-1234-123456789abc",
    "audioChunk": "BASE64_ENCODED_AUDIO_DATA"
  }'
```

#### 3. Close the streaming session

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