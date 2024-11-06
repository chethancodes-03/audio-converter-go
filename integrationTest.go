package integration

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

// Helper function to hash file contents
func hashFileContents(filePath string) (string, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:]), nil
}

// Test function for integration testing
func TestWAVtoFLACConversion(t *testing.T) {
	// Define input WAV files and expected FLAC files
	testCases := []struct {
		inputFilePath    string
		expectedFilePath string
	}{
		{"./samples/input1.wav", "./samples/expected1.flac"}, //given input three sample files along with three expected otput files
		{"./samples/input2.wav", "./samples/expected2.flac"},
		{"./samples/input3.wav", "./samples/expected3.flac"},
	}

	// Set up WebSocket connection to the server
	serverURL := "ws://localhost:3001/ws" // Replace with server address if different
	for _, tc := range testCases {
		t.Run(tc.inputFilePath, func(t *testing.T) {
			conn, _, err := websocket.Dial(context.Background(), serverURL, nil)
			if err != nil {
				t.Fatalf("failed to connect to WebSocket server: %v", err)
			}
			defer conn.Close(websocket.StatusInternalError, "test complete")

			// Open the WAV file
			wavFile, err := os.Open(tc.inputFilePath)
			if err != nil {
				t.Fatalf("failed to open input WAV file %s: %v", tc.inputFilePath, err)
			}
			defer wavFile.Close()

			// Read WAV file data and send over WebSocket
			wavData, err := ioutil.ReadAll(wavFile)
			if err != nil {
				t.Fatalf("failed to read input WAV file %s: %v", tc.inputFilePath, err)
			}

			// Send WAV data to the server in chunks
			err = wsjson.Write(context.Background(), conn, wavData)
			if err != nil {
				t.Fatalf("failed to send WAV data to server: %v", err)
			}

			// Receive FLAC data from the server
			var flacData bytes.Buffer
			for {
				_, message, err := conn.Read(context.Background())
				if err != nil {
					log.Println("Error reading from WebSocket:", err)
					break
				}
				flacData.Write(message)
			}

			// Save the received FLAC data to a temporary file
			outputFilePath := "./samples/output_temp.flac"
			err = ioutil.WriteFile(outputFilePath, flacData.Bytes(), 0644)
			if err != nil {
				t.Fatalf("failed to write received FLAC data to file: %v", err)
			}
			defer os.Remove(outputFilePath) // Clean up temporary file after test

			// Hash both the expected and actual output FLAC files
			expectedHash, err := hashFileContents(tc.expectedFilePath)
			if err != nil {
				t.Fatalf("failed to hash expected FLAC file %s: %v", tc.expectedFilePath, err)
			}

			actualHash, err := hashFileContents(outputFilePath)
			if err != nil {
				t.Fatalf("failed to hash output FLAC file: %v", err)
			}

			// Compare the hashes
			if expectedHash != actualHash {
				t.Errorf("FLAC conversion output does not match expected output for %s", tc.inputFilePath)
			}
		})
	}
}
