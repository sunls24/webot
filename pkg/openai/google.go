package openai

import (
	"fmt"
	"github.com/tidwall/gjson"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"webot/pkg/client"
)

const (
	nameGoogle = "GoogleSearch"
)

var Google = Tool{
	Name: nameGoogle,
	Func: fmt.Sprintf(`{
  "type": "function",
  "function": {
    "name": "%s",
    "strict": true,
    "description": "Using Google to search the internet",
    "parameters": {
      "type": "object",
      "properties": {
        "keyword": { "type": "string", "description": "Search keyword" }
      },
      "required": ["keyword"]
    }
  }
}`, nameGoogle),
	Call: callGoogle,
}

func callGoogle(keyword string) string {
	keyword = gjson.Get(keyword, "keyword").String()
	slog.Info("call google search", slog.String("keyword", keyword))

	var apiKey, _ = os.LookupEnv("GOOGLE_API_KEY")
	var engineId, _ = os.LookupEnv("GOOGLE_ENGINE_ID")
	if apiKey == "" || engineId == "" {
		slog.Warn("not found env GOOGLE_API_KEY or GOOGLE_ENGINE_ID")
		return "nothing"
	}
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("https://www.googleapis.com/customsearch/v1?&fields=items(title,link,snippet,pagemap/metatags(og:description))&key=%s&cx=%s&q=%s", apiKey, engineId, url.QueryEscape(keyword)), nil)
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("google search failed", slog.Any("err", err))
		return "nothing"
	}

	return gjson.GetBytes(resp, "items").String()
}
