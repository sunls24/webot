package openai

import (
	"encoding/json"
	"github.com/tidwall/gjson"
)

type Role string

const (
	RUser      Role = "user"
	RSystem    Role = "system"
	RAssistant Role = "assistant"
	RTool      Role = "tool"
)

type Message struct {
	Role    Role   `json:"role"`
	Content string `json:"content"`

	ToolCalls  json.RawMessage `json:"tool_calls,omitempty"`
	ToolCallId string          `json:"tool_call_id,omitempty"`
}

func NewMessage(msg gjson.Result) *Message {
	return &Message{
		Role:    Role(msg.Get("role").String()),
		Content: msg.Get("content").String(),
	}
}

type RequestBody struct {
	Model       string          `json:"model"`
	Messages    []Message       `json:"messages"`
	Temperature float64         `json:"temperature"`
	Tools       json.RawMessage `json:"tools,omitempty"`
}

type ToolCall func(string) string

type Tool struct {
	Name string
	Func string
	Call ToolCall
}
