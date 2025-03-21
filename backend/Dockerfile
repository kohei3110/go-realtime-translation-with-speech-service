FROM --platform=$BUILDPLATFORM golang:1.22-bullseye AS builder

# Install required dependencies
RUN dpkg --add-architecture amd64 && \
    apt-get update && apt-get install -y \
    build-essential \
    curl \
    pkg-config \
    crossbuild-essential-amd64 \
    uuid-dev:amd64 \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /build

# Copy and install Azure Speech SDK from local file
COPY speechsdk.tar.gz .
RUN tar -xzf speechsdk.tar.gz \
    && mkdir -p /usr/local/SpeechSDK \
    && cp -r SpeechSDK-Linux-*/include /usr/local/SpeechSDK/ \
    && cp SpeechSDK-Linux-*/lib/x64/libMicrosoft.CognitiveServices.Speech.core.so /usr/local/lib/ \
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
ENV CC=x86_64-linux-gnu-gcc
ENV CGO_CFLAGS="-I/usr/local/SpeechSDK/include/c_api"
ENV CGO_LDFLAGS="-L/usr/local/lib -lMicrosoft.CognitiveServices.Speech.core -luuid"
ENV LD_LIBRARY_PATH="/usr/local/lib:$LD_LIBRARY_PATH"

# Build the application
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o app ./cmd/api/main.go \
    && chmod +x app

# Start a new stage with explicit platform
FROM --platform=linux/amd64 debian:bullseye-slim

WORKDIR /app

# Copy the Speech SDK library from builder
COPY --from=builder /usr/local/lib/libMicrosoft.CognitiveServices.Speech.core.so /usr/local/lib/

# Install runtime dependencies with specific platform
RUN apt-get update && apt-get install -y \
    ca-certificates \
    uuid-runtime \
    libcurl4 \
    openssl \
    libasound2 \
    && rm -rf /var/lib/apt/lists/* \
    && ldconfig

# Copy the binary from builder with explicit destination
COPY --from=builder /build/app /app/app
RUN chmod +x /app/app

# Set default environment variables for application
ENV PORT=8080
ENV AZURE_SPEECH_KEY=""
ENV AZURE_SPEECH_REGION=""
ENV LD_LIBRARY_PATH="/usr/local/lib:$LD_LIBRARY_PATH"

# Additional environment variables for Speech SDK
ENV SPEECH_SDK_TRACE_LOG=1
ENV SPEECH_SDK_LOGGING_ENABLED=1

# Command to run the application with full path
CMD ["/app/app"]