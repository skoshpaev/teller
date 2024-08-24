package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

const (
	port         = "8080"                                                         // Port to be used in tests
	jwtSecret    = "H1g$eCr3t!2S#cUr3T@256-bSecr3tIt"                             // JWT secret to be used in tests
	execPath     = "../../bin/Teller/main"                                        // Path to the executable
	subscribeURL = "http://localhost:" + port + "/subscribe?channel=test-channel" // URL for SSE subscription
	publishURL   = "http://localhost:" + port + "/publish"                        // URL for message publishing
	token        = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.H1KnRyz0-_3OVZJJH-AYAwIMqR-9n5Uz9r97omtnGTc"
)

func TestStartWithoutSecretKey(t *testing.T) {
	cmd := exec.Command(execPath, "--port="+port)
	err := cmd.Run()
	if err == nil {
		t.Fatal("Expected error when starting without secret key, but got none")
	}

	defer cmd.Process.Kill()
}

func TestStartWithPortAndSecretKey(t *testing.T) {
	cmd := exec.Command(execPath, "--port="+port, "--jwt-secret="+jwtSecret)
	err := cmd.Start()
	if err != nil {
		t.Fatalf("Failed to start with port and secret key: %v", err)
	}

	defer cmd.Process.Kill()
}

func TestStartWithDuplicatePort(t *testing.T) {
	// Configure the first command with output to the console
	cmd1 := exec.Command(execPath, "--port="+port, "--jwt-secret="+jwtSecret)
	cmd1.Stdout = os.Stdout
	cmd1.Stderr = os.Stderr

	// Start the first instance
	err := cmd1.Start()
	if err != nil {
		t.Fatalf("Failed to start first instance: %v", err)
	}

	defer cmd1.Process.Kill()

	// Wait until the port is captured by the first instance
	for i := 0; i < 10; i++ { // Retry up to 10 times
		conn, err := net.Dial("tcp", "localhost:"+port)
		if err == nil {
			conn.Close()
			time.Sleep(1 * time.Second) // Wait a bit more to ensure the port is fully captured
			break
		}
		if i == 9 {
			t.Fatalf("First instance did not capture the port within the expected time: %v", err)
		}
		time.Sleep(1 * time.Second)
	}

	// Configure the second command with output to the console
	cmd2 := exec.Command(execPath, "--port="+port, "--jwt-secret="+jwtSecret)
	cmd2.Stdout = os.Stdout
	cmd2.Stderr = os.Stderr

	// Attempt to start the second instance
	err2 := cmd2.Run()

	if exitError, ok := err2.(*exec.ExitError); ok {
		code := exitError.ExitCode()
		t.Logf("Process exited with code: %d", code)
		if code == 0 {
			t.Fatal("Expected error when starting second instance on the same port, but got none")
		}
	} else {
		t.Logf("Second instance failed as expected: %v", err2)
	}

	defer cmd2.Process.Kill()
}

// ----------- Integration --------------------

type Message struct {
	Key string `json:"key"`
}

type PublishPayload struct {
	Channel string  `json:"channel"`
	Message Message `json:"message"`
}

func TestSSEIntegration(t *testing.T) {
	cmd := exec.Command(execPath, "--port="+port, "--jwt-secret="+jwtSecret)
	cmdErr := cmd.Start()
	if cmdErr != nil {
		t.Fatalf("Failed to start with port and secret key: %v", cmdErr)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a channel to receive SSE messages
	messageChannel := make(chan Message, 10)

	waitErr := waitForServerStart("http://localhost:"+port+"/", 20*time.Second)
	if waitErr != nil {
		t.Fatalf("Server did not start: %v", waitErr)
	}

	// Start a goroutine for SSE subscription
	go func() {
		err := subscribeToSSE(ctx, messageChannel)
		if err != nil {
			t.Errorf("Error subscribing to SSE: %v", err)
		}
	}()

	// Wait a bit to allow subscription to establish
	time.Sleep(10 * time.Second)

	// Send a message via HTTP publish
	expectedMessage := Message{Key: "value2"}
	err := publishMessage(expectedMessage)
	if err != nil {
		t.Errorf("Error publishing message via HTTP: %v", err)
	}

	// Check that the message was received via SSE
	select {
	case msg := <-messageChannel:
		t.Logf("Got: %+v", msg)
		if msg.Key != expectedMessage.Key {
			t.Errorf("Expected key '%s', but got '%s'", expectedMessage.Key, msg.Key)
		}
	case <-time.After(5 * time.Second):
		t.Error("Timeout waiting for message via SSE")
	}

	t.Logf("Closing main")

	_ = cmd.Process.Kill()
}

// Function to subscribe to SSE
func subscribeToSSE(ctx context.Context, messageChannel chan Message) error {
	req, err := http.NewRequest("GET", subscribeURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			line, err := reader.ReadBytes('\n')
			if err != nil {
				if err == io.EOF {
					return nil
				}
				return err
			}

			// Parse received line as SSE message
			if strings.HasPrefix(string(line), "data:") {
				var msg Message
				data := bytes.TrimPrefix(line, []byte("data:"))
				if err := json.Unmarshal(data, &msg); err != nil {
					return err
				}
				messageChannel <- msg
			}
		}
	}
}

// Function to publish a message via HTTP
func publishMessage(msg Message) error {
	payload := PublishPayload{
		Channel: "test-channel",
		Message: msg,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", publishURL, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Failed to send message, status: %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}

// Function to wait for server to start
func waitForServerStart(url string, timeout time.Duration) error {
	start := time.Now()
	for {
		if time.Since(start) > timeout {
			return fmt.Errorf("Server did not start within the allotted time")
		}

		resp, err := http.Get(url)
		if err == nil && resp.StatusCode == http.StatusNotFound {
			return nil // Consider server started if we get 404
		}

		time.Sleep(500 * time.Millisecond)
	}
}
