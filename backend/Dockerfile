FROM golang:1.22-bullseye

# Install required dependencies
RUN apt-get update && apt-get install -y \
    build-essential \
    curl \
    pkg-config \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy and install Azure Speech SDK from local file
COPY speechsdk.tar.gz .
RUN tar -xzf speechsdk.tar.gz \
    && mkdir -p /usr/local/SpeechSDK \
    && cp -r SpeechSDK-Linux-*/include /usr/local/SpeechSDK/ \
    && cp SpeechSDK-Linux-*/lib/arm64/libMicrosoft.CognitiveServices.Speech.core.so /usr/local/lib/ \
    && ldconfig \
    && rm -rf speechsdk.tar.gz SpeechSDK-Linux-*

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Set environment variables for compilation
ENV CGO_ENABLED=1
ENV CGO_CFLAGS="-I/usr/local/SpeechSDK/include/c_api"
ENV CGO_LDFLAGS="-L/usr/local/lib -lMicrosoft.CognitiveServices.Speech.core"
ENV LD_LIBRARY_PATH="/usr/local/lib:$LD_LIBRARY_PATH"

# Set default environment variables for application
ENV PORT=8080
ENV AZURE_SPEECH_KEY=""
ENV AZURE_SPEECH_REGION=""

# Build the application
RUN go build -o main ./cmd/api/main.go

# Command to run the application
CMD ["./main"]