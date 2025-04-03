package nlp

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Response struct {
	Intent string `json:"intent"`
	Error  string `json:"error,omitempty"`
}

var ErrNLP = errors.New("nlp service error")

// DetectIntentHTTP — отправляет запрос на NLP-сервис и возвращает интент.
func DetectIntentHTTP(nlpURL, text string) (string, error) {
	reqBody := map[string]string{"text": text}
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal error: %w", err)
	}

	resp, err := http.Post(nlpURL, "application/json", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("post to NLP failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response error: %w", err)
	}

	var nlpResp Response
	if err := json.Unmarshal(respBody, &nlpResp); err != nil {
		return "", fmt.Errorf("unmarshal error: %w", err)
	}

	if nlpResp.Error != "" {
		return "", fmt.Errorf("%w: %s", ErrNLP, nlpResp.Error)
	}

	intent := strings.ToLower(nlpResp.Intent)
	intent = strings.TrimSpace(intent)
	if intent == "" {
		return "unknown", nil
	}

	return intent, nil
}
