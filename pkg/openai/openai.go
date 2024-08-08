package openai

import (
	"fmt"
	"github.com/tidwall/gjson"
	"net/http"
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

func (oai *OpenAI) Chat(model string, messages []Message) (*Message, error) {
	const api = "/v1/chat/completions"

	reqBody := RequestBody{
		Model:       model,
		Messages:    messages,
		Temperature: 0.5,
	}

	req, _ := http.NewRequest(http.MethodPost, oai.opts.BaseURL+api, marshalBody(reqBody))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", oai.opts.APIKey))
	resp, err := do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", err, string(resp))
	}
	return NewMessage(gjson.GetBytes(resp, "choices.0.message")), nil
}
