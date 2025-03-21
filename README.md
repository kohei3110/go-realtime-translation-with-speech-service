# Real-time Translation Service

Backend API for real-time translation service using Azure Speech Service.

## Requirements

- Go 1.19 or higher
- Azure Speech Service account
- C++ compiler (gcc or clang)
- Azure Speech SDK for C/C++

### Installing Dependencies on macOS

```bash
# Install C++ compiler using Homebrew
brew install gcc

# Install Azure Speech SDK for C/C++
curl -L https://aka.ms/csspeech/macosbinary -o speechsdk.tar.gz
tar -xzf speechsdk.tar.gz
sudo mkdir -p /usr/local/include
sudo cp SpeechSDK-macOS/include/* /usr/local/include/
sudo cp SpeechSDK-macOS/lib/libMicrosoft.CognitiveServices.Speech.core.dylib /usr/local/lib/
rm -rf speechsdk.tar.gz SpeechSDK-macOS
```

## Setting Environment Variables

Set the following environment variables:

```bash
export PORT=8080  # API server port (optional, default: 8080)
export AZURE_SPEECH_KEY=your_key_here  # Azure Speech Service key
export AZURE_SPEECH_REGION=your_region_here  # Azure Speech Service region

# Set Azure Speech SDK library paths
export CGO_CFLAGS="-I/usr/local/include"
export CGO_LDFLAGS="-L/usr/local/lib -lMicrosoft.CognitiveServices.Speech.core"
export DYLD_LIBRARY_PATH="/usr/local/lib:$DYLD_LIBRARY_PATH"
```

## Setup Instructions

1. Clone the repository
```bash
git clone [repository-url]
cd go-realtime-translation-with-speech-service
```

2. Install dependencies
```bash
cd backend
go mod tidy
```

## Running with Docker

1. Create a .env file
```bash
AZURE_SPEECH_KEY=your_key_here
AZURE_SPEECH_REGION=your_region_here
```

2. Build and start Docker container
```bash
docker compose up --build
```

To stop the container, run the following command:
```bash
docker compose down
```

## Starting the API (Local Environment)

1. Navigate to the backend directory
```bash
cd backend  # if you're not already in the backend directory
```

2. Start the API server
```bash
go run cmd/api/main.go
```

When the server starts successfully, you will see a message like:
```
Server is running on port 8080
```

## Stopping the API

To stop the server, press `Ctrl+C` in the terminal. A graceful shutdown will be performed.

## Specifications

### Backend Specification

#### Technology Stack
- Language: Go 1.19+
- Framework: Standard library + Gorilla WebSocket
- External Service: Azure Speech Service
- Infrastructure: Docker

#### API Endpoints

1. **WebSocket Connection Endpoint**
   - Path: `/ws`
   - Method: GET (WebSocket Upgrade)
   - Function: Bidirectional communication for audio stream transmission and real-time translation

2. **Health Check Endpoint**
   - Path: `/health`
   - Method: GET
   - Response: `{"status": "ok"}`
   - Function: Service health verification

#### WebSocket Message Format

**From Client to Server:**
```json
{
  "type": "start_translation",
  "sourceLanguage": "ja-JP",
  "targetLanguage": "en-US",
  "audioFormat": "audio/wav"
}
```

```json
{
  "type": "audio_data",
  "data": "base64 encoded audio data"
}
```

```json
{
  "type": "stop_translation"
}
```

**From Server to Client:**
```json
{
  "type": "translation_result",
  "sourceText": "こんにちは",
  "translatedText": "Hello",
  "isFinal": true
}
```

```json
{
  "type": "error",
  "message": "error message"
}
```

#### Error Handling
- All errors are logged
- Error messages are returned to clients in JSON format
- Automatic reconnection attempts in case of connection errors

#### Performance Requirements
- Maximum concurrent connections: 100
- Latency: Within 1 second from voice input to translation display
- CPU usage: Average below 60%
- Memory usage: Maximum 512MB

### Frontend Specification

#### Technology Stack
- Language: TypeScript
- Framework: React
- Styling: CSS Modules or Tailwind CSS
- Build Tool: Vite

#### Functional Requirements

1. **User Interface**
   - Simple and intuitive UI
   - Responsive design (supporting mobile, tablet, desktop)
   - Dark mode/light mode toggle

2. **Voice Input**
   - Recording and sending microphone audio
   - Audio level indicator display
   - Automatic pause with silence detection

3. **Translation Display**
   - Simultaneous display of source language text and translated text
   - Saving and displaying translation history
   - Text copy function

4. **Settings**
   - Language pair selection (source and target languages)
   - Voice input sensitivity adjustment
   - Font size adjustment

5. **Status Display**
   - Connection status indicator
   - Error message display
   - Voice recognition status display

#### Non-functional Requirements
- Initial loading time: Within 2 seconds
- Offline functionality: Basic UI display and error messages
- Accessibility: WCAG 2.1 AA level compliance
- Mobile device battery consumption optimization

#### User Flow
1. Access the application
2. Select language pair
3. Grant microphone access permission
4. Click start button
5. Start speaking
6. Check translation results in real-time
7. Stop/resume as needed
8. Review or export translation history

#### Design Requirements
- Modern and clean interface
- Visual feedback provision
- Color contrast ratio: 4.5:1 or higher
- Icon and action button size: Minimum 44px×44px (for touch devices)