package termclient

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"log"
	"log/slog"
	"regexp"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/transform"
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
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // Allow self-signed certificates
		},
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
func (p *Client) sendKey(text string) error {
	// PTT WebSocket expects binary frames with specific line endings
	data := []byte(text)
	err := p.conn.WriteMessage(websocket.BinaryMessage, data)
	if err != nil {
		return err
	}
	log.Printf("Sent Key: %s", text)

	// Wait a moment for the server to process
	time.Sleep(time.Millisecond * 10)
	return nil
}

// sendByte sends a byte slice to PTT via WebSocket
func (p *Client) sendByte(b []byte) error {
	// PTT WebSocket expects binary frames with specific line endings
	err := p.conn.WriteMessage(websocket.BinaryMessage, b)
	if err != nil {
		return err
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

// readMessage reads and processes message from PTT WebSocket
func (p *Client) readMessage() (string, error) {
	slog.Info("Reading message from PTT")
	// Set a short read timeout
	p.conn.SetReadDeadline(time.Now().Add(10 * time.Second))

	_, rawData, err := p.conn.ReadMessage()
	if err != nil {
		slog.Warn("Error reading message", "err", err)
		return "", err
	}

	// Try to convert from Big5 to UTF-8
	big5Decoder := traditionalchinese.Big5.NewDecoder()
	utf8Data, _, err := transform.Bytes(big5Decoder, rawData)
	if err != nil {
		// If Big5 conversion fails, use raw data
		utf8Data = rawData
	}

	message := string(utf8Data)

	// Clean the message
	cleanedMessage := p.cleanMessage(message)
	return cleanedMessage, nil
}

// cleanMessage removes ANSI escape sequences and control characters
func (p *Client) cleanMessage(message string) string {
	// More comprehensive ANSI escape sequence removal
	// Remove standard ANSI escape sequences
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
	message = ansiRegex.ReplaceAllString(message, "")

	// Remove other ANSI sequences
	ansiRegex2 := regexp.MustCompile(`\x1b\([AB]`)
	message = ansiRegex2.ReplaceAllString(message, "")

	// Remove cursor movement and other control sequences
	ansiRegex3 := regexp.MustCompile(`\x1b\[[HfABCDsuJK]`)
	message = ansiRegex3.ReplaceAllString(message, "")

	// Clean up the message more carefully
	var buffer bytes.Buffer
	for i, r := range message {
		// Keep printable ASCII characters
		if r >= 32 && r <= 126 {
			buffer.WriteRune(r)
		} else if r == '\n' || r == '\r' || r == '\t' {
			// Keep basic whitespace
			buffer.WriteRune(r)
		} else if r > 127 {
			// For non-ASCII characters, check if they might be valid Chinese characters
			// This is a basic approach - in reality PTT uses Big5 encoding
			if r >= 0x4E00 && r <= 0x9FFF { // CJK Unified Ideographs range
				buffer.WriteRune(r)
			} else if r >= 0x3400 && r <= 0x4DBF { // CJK Extension A
				buffer.WriteRune(r)
			} else {
				// Replace other non-printable characters with space
				if i > 0 && i < len(message)-1 {
					buffer.WriteRune(' ')
				}
			}
		}
	}

	result := buffer.String()

	// Clean up multiple spaces and trim
	result = regexp.MustCompile(`\s+`).ReplaceAllString(result, " ")
	result = strings.TrimSpace(result)

	return result
}
