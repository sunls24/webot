package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

var client = &http.Client{
	Timeout: time.Minute,
}

func do(req *http.Request) ([]byte, error) {
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		return respBody, fmt.Errorf("bad status: %s", resp.Status)
	}

	return respBody, nil
}

func marshalBody(body any) io.Reader {
	data, _ := json.Marshal(body)
	return bytes.NewReader(data)
}
