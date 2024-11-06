# Audio Converter GO

## Overview

**Audio Converter GO** is a lightweight, Go-based audio file conversion tool utilizing [SoX (Sound eXchange)](http://sox.sourceforge.net/) to convert `.WAV` files to `.FLAC` format. It supports real-time streaming over WebSocket connections, enabling efficient handling of audio data for live conversions or large file streams.

This tool is designed to be flexible and easy to extend with additional audio processing features or formats.

---

## Features

- **WAV to FLAC Conversion**: Converts `.WAV` files to `.FLAC` format, leveraging SoX for high-quality processing.
- **Real-Time Streaming**: Streams audio data to handle large audio files or real-time conversions.
- **WebSocket API**: Uses WebSocket connections for audio streaming and conversion requests.
- **Concurrency Control**: Implements a semaphore to limit concurrent WebSocket connections for optimized resource usage.
- **Testing Suite**: Includes test files to validate functionality and ensure reliability.

---

## Technology Stack

- **Go**: Core language for server implementation, known for its performance and concurrency support.
- **SoX**: High-quality audio processing using SoX, a command-line utility for format conversion.
- **WebSocket**: Supports low-latency streaming of audio data, enabling continuous, bi-directional communication.

---

## Folder Structure

```plaintext
AUDIO CONVERTER-GO/
├── .venv/                   # Virtual environment (optional, for Python-based tests or scripts)
├── __debug_bin...exe        # Debug binaries created during VS Code sessions
├── file_example_WAV_2MG.wav # Sample input file for testing audio conversion
├── go.mod                   # Defines Go module and dependencies
├── go.sum                   # Contains dependency checksums
├── main.go                  # Main Go application code for the audio converter
├── maintest.go              # Go test file for testing functions in main.go
├── output.flac              # Converted output file in FLAC format
├── test_tone_440hz.wav      # Test file for frequency-based testing
├── test.py                  # Auxiliary Python script for testing or automation


Prerequisites
Go: Install Go, and verify with go version.
SoX (Sound eXchange): Install SoX, and verify with sox --version.
Git: Install Git if you intend to clone the repository.
WebSocket Client: Use tools like Postman or a custom WebSocket client script for testing.
Go Libraries and Dependencies
"github.com/gorilla/websocket": Manages WebSocket connections for real-time data streaming.
SoX Integration: SoX runs via command-line execution within the Go code; ensure it’s globally accessible.
Install Go dependencies with:

bash
Copy code
go mod tidy
API Endpoints and WebSocket Communication
WebSocket Endpoint
The application exposes a WebSocket endpoint for audio conversion and streaming:

Endpoint: ws://localhost:8080/convert
Method: WebSocket connection (Upgrade HTTP to WebSocket)
Input Format: .WAV audio data streamed as binary
Output Format: .FLAC audio data streamed as binary
WebSocket Workflow
Establish Connection: Connect to ws://localhost:8080/convert via WebSocket.
Send WAV Data: Stream .WAV audio data over the WebSocket.
Receive FLAC Data: The server converts .WAV to .FLAC in real-time and streams it back.
Close Connection: Disconnect when the conversion is complete to free resources.
WebSocket Message Protocol
Connection: Client initiates a WebSocket connection.
Binary Data Streaming: Audio data is sent as binary for efficient, continuous streaming.
Concurrency: Controlled via a semaphore that limits active WebSocket connections.
Error Handling: If an error occurs (e.g., unsupported format), the server sends an error message and closes the connection.
Setup Instructions
1. Clone the Repository
bash
Copy code
git clone https://github.com/chethancodes-03/audio-converter-go.git
cd audio-converter-go
2. Install Dependencies
bash
Copy code
go mod tidy
3. Run the WebSocket Server
Start the WebSocket server to enable streaming and conversion:

bash
Copy code
go run main.go
4. Testing Conversion
Use a WebSocket client to connect and stream .WAV data to ws://localhost:8080/convert.

Usage
Running Conversion
Start the Server: Run go run main.go to start the WebSocket server.
Stream Audio: Send .WAV files through WebSocket to receive real-time .FLAC conversions.
Conversion Output: Files are saved locally as .flac or streamed back to the client.
SoX Details
SoX is the core tool behind .WAV to .FLAC conversion. Key features include:

Sample Rate Adjustments: Change bit rate for various quality levels or storage efficiency.
Audio Effects: Options to add filters, equalizers, or noise suppression.
Error Handling: Automatically adjusts for inconsistencies, providing smooth conversions.
Modify the SoX command in main.go for custom parameters.

Testing
This project includes a basic testing suite to validate functionality:

bash
Copy code
go test maintest.go
Contributing
Fork the repository and submit pull requests for new features, bug fixes, or improvements.

Guidelines
Ensure existing functionality remains intact.
Add tests for new features.
Follow the project’s code style for consistency.
License
Licensed under the MIT License. See LICENSE for details.

Acknowledgments
SoX: Core audio conversion utility.
Gorilla WebSocket: WebSocket support for Go.
Go: Efficient concurrency and performance for server functionality.

