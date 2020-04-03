package message

// Message represents a message sent on the websocket server
type Message struct {
	Type    string                 `json:"type"`
	Payload map[string]interface{} `json:"payload"`
}
