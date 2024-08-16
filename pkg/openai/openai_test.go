package openai

import (
	"os"
	"testing"
)

var (
	baseURL, _ = os.LookupEnv("OPENAI_BASE_URL")
	apiKey, _  = os.LookupEnv("OPENAI_API_KEY")
	model, _   = os.LookupEnv("OPENAI_MODEL")
)

var openai = New(Options{
	BaseURL: baseURL,
	APIKey:  apiKey,
})

func TestOpenAI_Chat(t *testing.T) {
	if baseURL == "" || apiKey == "" {
		t.Fatal("openai base url or api key is empty")
	}
	msg, err := openai.Chat(model, []Message{
		{
			Role:    RUser,
			Content: "help me search linux.do",
		},
	}, Google)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s: %s\n", msg.Role, msg.Content)
}
