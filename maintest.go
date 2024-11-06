package main

import (
	"bytes"
	"errors"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"os/exec"
	"testing"
)

// Mock the exec.Command function to simulate SoX behavior
type MockExecCommand struct {
	mock.Mock
}

func (m *MockExecCommand) Start() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockExecCommand) Wait() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockExecCommand) StdoutPipe() (io.Reader, error) {
	args := m.Called()
	return args.Get(0).(io.Reader), args.Error(1)
}

func (m *MockExecCommand) StdinPipe() (io.Writer, error) {
	args := m.Called()
	return args.Get(0).(io.Writer), args.Error(1)
}

// Unit test to test the conversion logic
func TestWavToFlacConversion(t *testing.T) {
	// Create a mock of the exec.Command
	mockCmd := new(MockExecCommand)

	// Simulate SoX starting successfully
	mockCmd.On("Start").Return(nil)

	// Simulate the SoX command output (a byte slice with fake FLAC data)
	flacData := []byte("fake FLAC data")
	mockCmd.On("StdoutPipe").Return(bytes.NewReader(flacData), nil)
	mockCmd.On("StdinPipe").Return(new(bytes.Buffer), nil)

	// Simulate that SoX finishes without errors
	mockCmd.On("Wait").Return(nil)

	// Replace the exec.Command with our mock
	execCommand = func(name string, arg ...string) *exec.Cmd {
		return mockCmd
	}

	// Simulate the WebSocket behavior
	fakeWavData := []byte("fake WAV data") // Simulate incoming WAV data
	conn := &websocket.Conn{}

	// Function under test (simplified version of the WebSocket handler)
	err := handleWavToFlacConversion(fakeWavData, conn)

	// Check for any errors
	assert.Nil(t, err)

	// Validate that SoX was started and data was sent to it
	mockCmd.AssertExpectations(t)
}

// Unit test to test the WebSocket connection handling
func TestWebSocketConnectionHandling(t *testing.T) {
	// Create a mock WebSocket connection
	mockConn := new(MockWebSocketConn)

	// Simulate WebSocket reading and writing
	mockConn.On("ReadMessage").Return([]byte("fake WAV data"), nil)
	mockConn.On("WriteMessage", websocket.BinaryMessage, []byte("fake FLAC data")).Return(nil)

	// Simulate the WebSocket handler function
	err := handleWebSocket(mockConn)

	// Check that the connection was handled correctly
	assert.Nil(t, err)

	// Validate that the right data was sent and received
	mockConn.AssertExpectations(t)
}

// Handle the WAV to FLAC conversion logic (simplified version)
func handleWavToFlacConversion(wavData []byte, conn *websocket.Conn) error {
	// Create the SoX command (using the mock in the test)
	cmd := exec.Command("sox", "-t", "wav", "-", "-t", "flac", "-")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	// Start the SoX command
	if err := cmd.Start(); err != nil {
		return err
	}

	// Write the WAV data to SoX for conversion
	_, err = stdin.Write(wavData)
	if err != nil {
		return err
	}

	// Read the FLAC data from SoX output and send it back to the WebSocket
	flacData := make([]byte, 8092)
	_, err = stdout.Read(flacData)
	if err != nil {
		return err
	}

	// Send the FLAC data to the WebSocket
	return conn.WriteMessage(websocket.BinaryMessage, flacData)
}

// Mock WebSocket connection for testing
type MockWebSocketConn struct {
	mock.Mock
}

func (m *MockWebSocketConn) WriteMessage(messageType int, p []byte) error {
	args := m.Called(messageType, p)
	return args.Error(0)
}

func (m *MockWebSocketConn) ReadMessage() (messageType int, p []byte, err error) {
	args := m.Called()
	return args.Int(0), args.Get(1).([]byte), args.Error(2)
}

func (m *MockWebSocketConn) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockWebSocketConn) LocalAddr() net.Addr {
	args := m.Called()
	return args.Get(0).(net.Addr)
}

func (m *MockWebSocketConn) RemoteAddr() net.Addr {
	args := m.Called()
	return args.Get(0).(net.Addr)
}

func (m *MockWebSocketConn) SetDeadline(t time.Time) error {
	args := m.Called(t)
	return args.Error(0)
}

func (m *MockWebSocketConn) SetReadDeadline(t time.Time) error {
	args := m.Called(t)
	return args.Error(0)
}

func (m *MockWebSocketConn) SetWriteDeadline(t time.Time) error {
	args := m.Called(t)
	return args.Error(0)
}
