package openai

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"net/http"
	"strings"
	"webot/pkg/client"
)

type OpenAI struct {
	opts Options
}

type Options struct {
	BaseURL string
	APIKey  string
}

func New(opts Options) *OpenAI {
	return &OpenAI{
		opts: opts,
	}
}

func (oai *OpenAI) Chat(model string, messages []Message, tools ...Tool) (*Message, error) {
	const api = "/v1/chat/completions"

	var toolsBody json.RawMessage
	if len(tools) != 0 {
		var funcList = make([]string, 0, len(tools))
		for _, t := range tools {
			funcList = append(funcList, t.Func)
		}
		toolsBody = json.RawMessage(fmt.Sprintf("[%s]", strings.Join(funcList, ",")))
	}

	reqBody := RequestBody{
		Model:       model,
		Messages:    messages,
		Temperature: 0.5,
		Tools:       toolsBody,
	}

	req, _ := http.NewRequest(http.MethodPost, oai.opts.BaseURL+api, client.MarshalBody(reqBody))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", oai.opts.APIKey))
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	result := gjson.GetBytes(resp, "choices.0.message")
	var toolCalls = result.Get("tool_calls")
	var toolCallsResult []Message
	for _, call := range toolCalls.Array() {
		callId := call.Get("id").String()
		call = call.Get("function")
		name := call.Get("name").String()
		args := call.Get("arguments").String()

		for _, t := range tools {
			if t.Name != name {
				continue
			}
			toolCallsResult = append(toolCallsResult, Message{
				Role:       RTool,
				Content:    t.Call(args),
				ToolCallId: callId,
			})
		}
	}
	if len(toolCallsResult) != 0 {
		messages = append(messages, Message{
			Role:      RAssistant,
			ToolCalls: json.RawMessage(toolCalls.Raw),
		})
		messages = append(messages, toolCallsResult...)
		return oai.Chat(model, messages, tools...)
	}

	message := NewMessage(result)
	if message.Role == "" {
		return nil, errors.New(string(resp))
	}
	return message, nil
}
