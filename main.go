package main

import (
	"fmt"
	"io"
	"log"
	"os/exec"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

const bufferSize = 8092   // Buffer size for reading FLAC data chunks
const maxConnections = 10 // Max number of concurrent WebSocket connections

var semaphore = make(chan struct{}, maxConnections) // Semaphore to limit connections

func main() {
	app := fiber.New()

	// Middleware to validate request origin
	app.Use("/ws", func(c *fiber.Ctx) error {
		host := c.Get("host")
		if host == "localhost:3001" || host == "127.0.0.1:3001" {
			c.Locals("Host", host)
			return c.Next()
		}
		return c.Status(403).SendString("Request origin not allowed")
	})

	// WebSocket endpoint to handle real-time audio conversion
	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		// Attempt to acquire a connection slot
		select {
		case semaphore <- struct{}{}:
			defer func() { <-semaphore }() // Release slot on exit
		default:
			log.Println("Connection limit reached; rejecting new WebSocket connection")
			c.Close() // Close connection if limit is reached
			return
		}

		fmt.Println("WebSocket connection opened")
		defer fmt.Println("WebSocket connection closed")

		// Set up SoX for WAV to FLAC conversion
		cmd := exec.Command("C:\\Program Files (x86)\\sox-14-4-2\\sox.exe", "-t", "wav", "-", "-t", "flac", "-")
		soxOut, err := cmd.StdoutPipe()
		if err != nil {
			log.Println("Error setting up SoX output:", err)
			return
		}
		soxIn, err := cmd.StdinPipe()
		if err != nil {
			log.Println("Error setting up SoX input:", err)
			return
		}

		// Start the SoX process
		if err := cmd.Start(); err != nil {
			log.Println("Error starting SoX process:", err)
			return
		}
		defer cmd.Wait() // Ensure SoX completes before cleanup

		// Channel to signal when WebSocket connection is closed
		done := make(chan struct{})

		// Goroutine to read WAV data from WebSocket and write to SoX
		go func() {
			defer close(done) // Signal done when this goroutine exits
			for {
				// Read WAV data from WebSocket
				_, msg, err := c.ReadMessage()
				if err != nil {
					log.Println("Error reading from WebSocket:", err)
					return
				}

				// Write WAV data to SoX's input for conversion
				if _, err := soxIn.Write(msg); err != nil {
					log.Println("Error writing to SoX:", err)
					return
				}
			}
		}()

		// Buffer for holding FLAC data chunks to send back to the client
		flacBuffer := make([]byte, bufferSize)
		for {
			select {
			case <-done:
				// WebSocket connection closed, clean up
				soxIn.Close() // Close SoX input to signal no more data is coming
				return
			default:
				// Read converted FLAC data from SoX output
				n, err := soxOut.Read(flacBuffer)
				if err != nil {
					if err == io.EOF {
						log.Println("SoX conversion completed.")
					} else {
						log.Println("Error reading from SoX output:", err)
					}
					return
				}

				// Send FLAC data back to WebSocket in real-time
				if err := c.WriteMessage(websocket.BinaryMessage, flacBuffer[:n]); err != nil {
					log.Println("Error writing to WebSocket:", err)
					return
				}
			}
		}
	}))

	// Start the Fiber server on localhost at port 3001
	log.Fatal(app.Listen(":3001"))
}
