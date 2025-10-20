package termclient

import (
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// PTTClient represents a PTT WebSocket client
type Client struct {
	conn *websocket.Conn
}

// NewPTTClient creates a new PTT client
func NewClient() *Client {
	return &Client{}
}

// Connect establishes WebSocket connection to PTT
func (p *Client) Connect() error {
	endpoint := "wss://ws.ptt.cc/bbs"

	log.Printf("Connecting to %s", endpoint)

	// Configure WebSocket dialer
	dialer := websocket.Dialer{
		HandshakeTimeout: 15 * time.Second,
	}

	// Set headers that PTT expects
	headers := map[string][]string{
		"Origin":     {"https://term.ptt.cc"},
		"User-Agent": {"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"},
	}

	// Attempt connection
	conn, response, err := dialer.Dial(endpoint, headers)
	if err != nil {
		log.Printf("Failed to connect to %s: %v", endpoint, err)
		if response != nil {
			log.Printf("Response status: %s", response.Status)
		}
		return fmt.Errorf("failed to connect to PTT WebSocket: %v", err)
	}
	conn.SetCloseHandler(func(code int, text string) error {
		log.Printf("WebSocket closed with code %d: %s", code, text)
		return nil
	})

	p.conn = conn
	log.Printf("Successfully connected to PTT WebSocket")
	return nil
}

func (p *Client) Disconnect() error {
	if p.conn != nil {
		err := p.conn.Close()
		if err != nil {
			return fmt.Errorf("failed to close WebSocket connection: %v", err)
		}
		log.Printf("WebSocket connection closed")
	}
	return nil
}

// sendText sends text to PTT via WebSocket
func (p *Client) sendText(text string) error {
	// PTT WebSocket expects binary frames with specific line endings
	data := []byte(text)
	err := p.conn.WriteMessage(websocket.BinaryMessage, data)
	if err != nil {
		return err
	}
	// log.Printf("Sent Text: %s", text)
	return nil
}

// sendText sends text to PTT via WebSocket
func (p *Client) sendTextBySequence(text string, delay time.Duration) error {
	// PTT WebSocket expects binary frames with specific line endings
	data := []byte(text)
	for _, b := range data {
		err := p.conn.WriteMessage(websocket.BinaryMessage, []byte{b})
		if err != nil {
			return err
		}
		time.Sleep(delay)
	}
	// log.Printf("Sent Text: %s", text)
	return nil
}
