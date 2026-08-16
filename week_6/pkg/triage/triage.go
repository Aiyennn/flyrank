package triage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type TriageResponse struct {
	Category    string   `json:"category"`
	Sentiment   string   `json:"sentiment"`
	Urgency     string   `json:"urgency"`
	ActionItems []string `json:"action_items"`
	Summary     string   `json:"summary"`
}

func (r *TriageResponse) Validate() error {
	category := strings.ToLower(strings.TrimSpace(r.Category))
	validCategories := map[string]bool{
		"billing":         true,
		"technical":       true,
		"account":         true,
		"feature_request": true,
		"spam":            true,
	}
	if !validCategories[category] {
		return fmt.Errorf("invalid category %q (must be one of: billing, technical, account, feature_request, spam)", r.Category)
	}

	sentiment := strings.ToLower(strings.TrimSpace(r.Sentiment))
	validSentiments := map[string]bool{
		"positive": true,
		"neutral":  true,
		"negative": true,
	}
	if !validSentiments[sentiment] {
		return fmt.Errorf("invalid sentiment %q (must be one of: positive, neutral, negative)", r.Sentiment)
	}

	urgency := strings.ToLower(strings.TrimSpace(r.Urgency))
	validUrgencies := map[string]bool{
		"low":      true,
		"medium":   true,
		"high":     true,
		"critical": true,
	}
	if !validUrgencies[urgency] {
		return fmt.Errorf("invalid urgency %q (must be one of: low, medium, high, critical)", r.Urgency)
	}

	if len(r.ActionItems) == 0 {
		return fmt.Errorf("action_items list cannot be empty")
	}
	for i, item := range r.ActionItems {
		if strings.TrimSpace(item) == "" {
			return fmt.Errorf("action item at index %d is empty", i)
		}
	}

	if strings.TrimSpace(r.Summary) == "" {
		return fmt.Errorf("summary cannot be empty")
	}

	return nil
}

type GeminiRequest struct {
	Contents         []Content         `json:"contents"`
	GenerationConfig *GenerationConfig `json:"generationConfig,omitempty"`
}

type Content struct {
	Parts []Part `json:"parts"`
}

type Part struct {
	Text string `json:"text"`
}

type GenerationConfig struct {
	ResponseMIMEType string  `json:"responseMimeType,omitempty"`
	ResponseSchema   *Schema `json:"responseSchema,omitempty"`
}

type Schema struct {
	Type        string             `json:"type"`
	Description string             `json:"description,omitempty"`
	Properties  map[string]*Schema `json:"properties,omitempty"`
	Required    []string           `json:"required,omitempty"`
	Items       *Schema            `json:"items,omitempty"`
}

type GeminiResponse struct {
	Candidates []Candidate `json:"candidates"`
}

type Candidate struct {
	Content      *Content `json:"content"`
	FinishReason string   `json:"finishReason"`
}

type LLMError struct {
	StatusCode int
	Message    string
	Retryable  bool
}

func (e *LLMError) Error() string {
	return fmt.Sprintf("LLM error (status %d): %s (retryable: %t)", e.StatusCode, e.Message, e.Retryable)
}

type Client struct {
	APIKey     string
	BaseURL    string
	Model      string
	HTTPClient *http.Client
	MaxRetries int
	MinDelay   time.Duration
	MaxDelay   time.Duration
	Timeout    time.Duration
}

func NewClient(apiKey string) *Client {
	return &Client{
		APIKey:  apiKey,
		BaseURL: "https://generativelanguage.googleapis.com",
		Model:   "gemini-3.7-flash",
		HTTPClient: &http.Client{
			Timeout: 15 * time.Second,
		},
		MaxRetries: 3,
		MinDelay:   500 * time.Millisecond,
		MaxDelay:   3000 * time.Millisecond,
		Timeout:    8 * time.Second,
	}
}

func (c *Client) calculateBackoff(attempt int) time.Duration {
	if attempt <= 0 {
		return c.MinDelay
	}
	delay := c.MinDelay * time.Duration(1<<attempt)
	if delay > c.MaxDelay {
		delay = c.MaxDelay
	}
	halfDelay := delay / 2
	jitter := rand.Int63n(int64(halfDelay))
	return halfDelay + time.Duration(jitter)
}

func (c *Client) TriageTicket(ctx context.Context, ticketText string) (*TriageResponse, error) {
	if c.APIKey == "" {
		return nil, &LLMError{StatusCode: 401, Message: "API Key is required but not provided", Retryable: false}
	}

	geminiSchema := &Schema{
		Type: "OBJECT",
		Properties: map[string]*Schema{
			"category": {
				Type:        "STRING",
				Description: "Category of the support ticket: billing, technical, account, feature_request, or spam",
			},
			"sentiment": {
				Type:        "STRING",
				Description: "User's sentiment: positive, neutral, or negative",
			},
			"urgency": {
				Type:        "STRING",
				Description: "Urgency level: low, medium, high, or critical",
			},
			"action_items": {
				Type: "ARRAY",
				Items: &Schema{
					Type: "STRING",
				},
				Description: "Actionable next steps to resolve the ticket",
			},
			"summary": {
				Type:        "STRING",
				Description: "A 1-2 sentence concise summary of the issue",
			},
		},
		Required: []string{"category", "sentiment", "urgency", "action_items", "summary"},
	}

	reqBody := GeminiRequest{
		Contents: []Content{
			{
				Parts: []Part{
					{
						Text: fmt.Sprintf("Triage the following support ticket:\n\n%s", ticketText),
					},
				},
			},
		},
		GenerationConfig: &GenerationConfig{
			ResponseMIMEType: "application/json",
			ResponseSchema:   geminiSchema,
		},
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	var lastErr error
	for attempt := 0; attempt <= c.MaxRetries; attempt++ {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		if attempt > 0 {
			delay := c.calculateBackoff(attempt - 1)
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(delay):
			}
		}

		attemptCtx, cancel := context.WithTimeout(ctx, c.Timeout)
		respBytes, err := c.doCall(attemptCtx, reqBytes)
		cancel()

		if err != nil {
			lastErr = err
			if llmErr, ok := err.(*LLMError); ok && !llmErr.Retryable {
				return nil, err
			}
			continue
		}

		triageResp, err := c.parseAndValidate(respBytes)
		if err != nil {
			lastErr = fmt.Errorf("parse/validation error on attempt %d: %w", attempt+1, err)
			continue
		}

		return triageResp, nil
	}

	return nil, fmt.Errorf("failed after %d attempts: %w", c.MaxRetries+1, lastErr)
}

func (c *Client) doCall(ctx context.Context, reqBytes []byte) ([]byte, error) {
	url := fmt.Sprintf("%s/v1beta/models/%s:generateContent?key=%s", c.BaseURL, c.Model, c.APIKey)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(reqBytes))
	if err != nil {
		return nil, &LLMError{StatusCode: 0, Message: fmt.Sprintf("failed to create request: %s", err.Error()), Retryable: false}
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		retryable := true
		if ctx.Err() == context.Canceled {
			retryable = false
		}
		return nil, &LLMError{StatusCode: 0, Message: err.Error(), Retryable: retryable}
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &LLMError{StatusCode: resp.StatusCode, Message: fmt.Sprintf("failed to read response: %s", err.Error()), Retryable: true}
	}

	if resp.StatusCode != http.StatusOK {
		retryable := false
		if resp.StatusCode == http.StatusTooManyRequests || (resp.StatusCode >= 500 && resp.StatusCode < 600) {
			retryable = true
		}
		return nil, &LLMError{
			StatusCode: resp.StatusCode,
			Message:    fmt.Sprintf("API returned status %s: %s", resp.Status, string(respBytes)),
			Retryable:  retryable,
		}
	}

	return respBytes, nil
}

func (c *Client) parseAndValidate(respBytes []byte) (*TriageResponse, error) {
	var geminiResp GeminiResponse
	if err := json.Unmarshal(respBytes, &geminiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Gemini response: %w", err)
	}

	if len(geminiResp.Candidates) == 0 {
		return nil, fmt.Errorf("no candidates returned in response")
	}

	candidate := geminiResp.Candidates[0]
	if candidate.Content == nil || len(candidate.Content.Parts) == 0 {
		return nil, fmt.Errorf("empty candidate content returned")
	}

	jsonText := candidate.Content.Parts[0].Text
	if jsonText == "" {
		return nil, fmt.Errorf("empty text part returned")
	}

	var triageResp TriageResponse
	if err := json.Unmarshal([]byte(strings.TrimSpace(jsonText)), &triageResp); err != nil {
		return nil, fmt.Errorf("failed to parse inner response JSON: %w (raw response: %s)", err, jsonText)
	}

	if err := triageResp.Validate(); err != nil {
		return nil, fmt.Errorf("schema validation failure: %w", err)
	}

	return &triageResp, nil
}

type TriageRequest struct {
	Text string `json:"text"`
}

type APIErrorResponse struct {
	Error string `json:"error"`
}

func (c *Client) HandleTriageTicket(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(APIErrorResponse{Error: "Only POST requests are allowed"})
		return
	}

	var reqBody TriageRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(APIErrorResponse{Error: "Invalid JSON request body"})
		return
	}

	if strings.TrimSpace(reqBody.Text) == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(APIErrorResponse{Error: "Field 'text' cannot be empty"})
		return
	}

	triageResult, err := c.TriageTicket(r.Context(), reqBody.Text)
	if err != nil {
		statusCode := http.StatusInternalServerError
		message := err.Error()

		if llmErr, ok := err.(*LLMError); ok {
			switch llmErr.StatusCode {
			case http.StatusUnauthorized:
				statusCode = http.StatusBadGateway
				message = "Upstream LLM authorization failed"
			case http.StatusTooManyRequests:
				statusCode = http.StatusTooManyRequests
				message = "Rate limit exceeded, try again later"
			}
		}

		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(APIErrorResponse{Error: message})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(triageResult)
}
