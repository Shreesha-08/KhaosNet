package server

type IncomingMessage struct {
	Command string   `json:"command,omitempty"`
	Args    []string `json:"args,omitempty"`
	Text    string   `json:"text,omitempty"`
	Room    string   `json:"room,omitempty"`
}

type OutgoingMessage struct {
	Type      string         `json:"type"`
	From      string         `json:"from,omitempty"`
	Room      string         `json:"room,omitempty"`
	Text      string         `json:"text,omitempty"`
	Data      map[string]any `json:"data,omitempty"`
	Timestamp int64          `json:"timestamp,omitempty"`
}
