package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

var Client = &http.Client{
	Timeout: time.Minute,
}

func Do(req *http.Request) ([]byte, error) {
	resp, err := Client.Do(req)
	if err != nil {
		return nil, err
	}
	respBody, err := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		status := fmt.Sprintf("bad status: %s", resp.Status)
		if len(respBody) != 0 {
			status = fmt.Sprintf("%s, body: %s", status, string(respBody))
		}
		return nil, errors.New(status)
	}

	return respBody, nil
}

func MarshalBody(body any) io.Reader {
	data, _ := json.Marshal(body)
	return bytes.NewReader(data)
}
