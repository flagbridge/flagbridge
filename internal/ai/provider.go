package ai

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Provider is the interface for LLM provider implementations.
type Provider interface {
	// Chat sends a non-streaming request and returns the full response.
	Chat(req ProviderRequest) (*ProviderResponse, error)
	// ChatStream sends a streaming request and returns a channel of chunks.
	ChatStream(req ProviderRequest) (<-chan StreamChunk, error)
	// Name returns the provider identifier.
	Name() string
}

// NewProvider creates the appropriate provider based on name.
func NewProvider(name, apiKey, baseURL string) (Provider, error) {
	switch name {
	case "anthropic":
		return newAnthropicProvider(apiKey, baseURL), nil
	case "openai":
		return newOpenAIProvider(apiKey, baseURL), nil
	case "ollama":
		return newOllamaProvider(baseURL), nil
	default:
		return nil, fmt.Errorf("unsupported provider: %s", name)
	}
}

// --- Anthropic ---

type anthropicProvider struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

func newAnthropicProvider(apiKey, baseURL string) *anthropicProvider {
	if baseURL == "" {
		baseURL = "https://api.anthropic.com"
	}
	return &anthropicProvider{apiKey: apiKey, baseURL: baseURL, client: &http.Client{}}
}

func (p *anthropicProvider) Name() string { return "anthropic" }

func (p *anthropicProvider) Chat(req ProviderRequest) (*ProviderResponse, error) {
	body := p.buildBody(req)
	httpReq, err := http.NewRequest("POST", p.baseURL+"/v1/messages", jsonBody(body))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	p.setHeaders(httpReq)

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("calling Anthropic API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("anthropic API error (status %d): %s", resp.StatusCode, string(b))
	}

	var result struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
		Usage struct {
			InputTokens  int `json:"input_tokens"`
			OutputTokens int `json:"output_tokens"`
		} `json:"usage"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	var text string
	for _, c := range result.Content {
		text += c.Text
	}

	return &ProviderResponse{
		Content: text,
		Usage:   &Usage{InputTokens: result.Usage.InputTokens, OutputTokens: result.Usage.OutputTokens},
	}, nil
}

func (p *anthropicProvider) ChatStream(req ProviderRequest) (<-chan StreamChunk, error) {
	body := p.buildBody(req)
	body["stream"] = true

	httpReq, err := http.NewRequest("POST", p.baseURL+"/v1/messages", jsonBody(body))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	p.setHeaders(httpReq)

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("calling Anthropic API: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("anthropic API error (status %d): %s", resp.StatusCode, string(b))
	}

	ch := make(chan StreamChunk, 64)
	go func() {
		defer resp.Body.Close()
		defer close(ch)
		p.readAnthropicStream(resp.Body, ch)
	}()

	return ch, nil
}

func (p *anthropicProvider) readAnthropicStream(r io.Reader, ch chan<- StreamChunk) {
	scanner := bufio.NewScanner(r)
	var usage *Usage

	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data: ") {
			continue
		}
		data := strings.TrimPrefix(line, "data: ")

		var event struct {
			Type  string `json:"type"`
			Delta struct {
				Text string `json:"text"`
			} `json:"delta"`
			Usage struct {
				InputTokens  int `json:"input_tokens"`
				OutputTokens int `json:"output_tokens"`
			} `json:"usage"`
		}
		if err := json.Unmarshal([]byte(data), &event); err != nil {
			continue
		}

		switch event.Type {
		case "content_block_delta":
			if event.Delta.Text != "" {
				ch <- StreamChunk{Content: event.Delta.Text}
			}
		case "message_delta":
			if event.Usage.OutputTokens > 0 {
				usage = &Usage{OutputTokens: event.Usage.OutputTokens}
			}
		case "message_stop":
			ch <- StreamChunk{Done: true, Usage: usage}
			return
		case "error":
			ch <- StreamChunk{Err: fmt.Errorf("anthropic stream error")}
			return
		}
	}
}

func (p *anthropicProvider) buildBody(req ProviderRequest) map[string]any {
	msgs := make([]map[string]string, len(req.Messages))
	for i, m := range req.Messages {
		msgs[i] = map[string]string{"role": m.Role, "content": m.Content}
	}
	body := map[string]any{
		"model":      req.Model,
		"messages":   msgs,
		"max_tokens": req.MaxTokens,
	}
	if req.System != "" {
		body["system"] = req.System
	}
	if req.Temperature > 0 {
		body["temperature"] = req.Temperature
	}
	return body
}

func (p *anthropicProvider) setHeaders(r *http.Request) {
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("x-api-key", p.apiKey)
	r.Header.Set("anthropic-version", "2023-06-01")
}

// --- OpenAI ---

type openAIProvider struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

func newOpenAIProvider(apiKey, baseURL string) *openAIProvider {
	if baseURL == "" {
		baseURL = "https://api.openai.com"
	}
	return &openAIProvider{apiKey: apiKey, baseURL: baseURL, client: &http.Client{}}
}

func (p *openAIProvider) Name() string { return "openai" }

func (p *openAIProvider) Chat(req ProviderRequest) (*ProviderResponse, error) {
	body := p.buildBody(req)
	httpReq, err := http.NewRequest("POST", p.baseURL+"/v1/chat/completions", jsonBody(body))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	p.setHeaders(httpReq)

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("calling OpenAI API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("openai API error (status %d): %s", resp.StatusCode, string(b))
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Usage struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
		} `json:"usage"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	var text string
	if len(result.Choices) > 0 {
		text = result.Choices[0].Message.Content
	}

	return &ProviderResponse{
		Content: text,
		Usage:   &Usage{InputTokens: result.Usage.PromptTokens, OutputTokens: result.Usage.CompletionTokens},
	}, nil
}

func (p *openAIProvider) ChatStream(req ProviderRequest) (<-chan StreamChunk, error) {
	body := p.buildBody(req)
	body["stream"] = true

	httpReq, err := http.NewRequest("POST", p.baseURL+"/v1/chat/completions", jsonBody(body))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	p.setHeaders(httpReq)

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("calling OpenAI API: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("openai API error (status %d): %s", resp.StatusCode, string(b))
	}

	ch := make(chan StreamChunk, 64)
	go func() {
		defer resp.Body.Close()
		defer close(ch)
		p.readOpenAIStream(resp.Body, ch)
	}()

	return ch, nil
}

func (p *openAIProvider) readOpenAIStream(r io.Reader, ch chan<- StreamChunk) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data: ") {
			continue
		}
		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			ch <- StreamChunk{Done: true}
			return
		}

		var event struct {
			Choices []struct {
				Delta struct {
					Content string `json:"content"`
				} `json:"delta"`
				FinishReason *string `json:"finish_reason"`
			} `json:"choices"`
			Usage *struct {
				PromptTokens     int `json:"prompt_tokens"`
				CompletionTokens int `json:"completion_tokens"`
			} `json:"usage"`
		}
		if err := json.Unmarshal([]byte(data), &event); err != nil {
			continue
		}

		if len(event.Choices) > 0 && event.Choices[0].Delta.Content != "" {
			ch <- StreamChunk{Content: event.Choices[0].Delta.Content}
		}
	}
}

func (p *openAIProvider) buildBody(req ProviderRequest) map[string]any {
	msgs := make([]map[string]string, 0, len(req.Messages)+1)
	if req.System != "" {
		msgs = append(msgs, map[string]string{"role": "system", "content": req.System})
	}
	for _, m := range req.Messages {
		msgs = append(msgs, map[string]string{"role": m.Role, "content": m.Content})
	}
	body := map[string]any{
		"model":      req.Model,
		"messages":   msgs,
		"max_tokens": req.MaxTokens,
	}
	if req.Temperature > 0 {
		body["temperature"] = req.Temperature
	}
	return body
}

func (p *openAIProvider) setHeaders(r *http.Request) {
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Authorization", "Bearer "+p.apiKey)
}

// --- Ollama ---

type ollamaProvider struct {
	baseURL string
	client  *http.Client
}

func newOllamaProvider(baseURL string) *ollamaProvider {
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}
	return &ollamaProvider{baseURL: baseURL, client: &http.Client{}}
}

func (p *ollamaProvider) Name() string { return "ollama" }

func (p *ollamaProvider) Chat(req ProviderRequest) (*ProviderResponse, error) {
	body := p.buildBody(req, false)
	httpReq, err := http.NewRequest("POST", p.baseURL+"/api/chat", jsonBody(body))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("calling Ollama API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ollama API error (status %d): %s", resp.StatusCode, string(b))
	}

	var result struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &ProviderResponse{Content: result.Message.Content}, nil
}

func (p *ollamaProvider) ChatStream(req ProviderRequest) (<-chan StreamChunk, error) {
	body := p.buildBody(req, true)
	httpReq, err := http.NewRequest("POST", p.baseURL+"/api/chat", jsonBody(body))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("calling Ollama API: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("ollama API error (status %d): %s", resp.StatusCode, string(b))
	}

	ch := make(chan StreamChunk, 64)
	go func() {
		defer resp.Body.Close()
		defer close(ch)

		decoder := json.NewDecoder(resp.Body)
		for {
			var chunk struct {
				Message struct {
					Content string `json:"content"`
				} `json:"message"`
				Done bool `json:"done"`
			}
			if err := decoder.Decode(&chunk); err != nil {
				if err != io.EOF {
					ch <- StreamChunk{Err: err}
				}
				return
			}

			if chunk.Message.Content != "" {
				ch <- StreamChunk{Content: chunk.Message.Content}
			}
			if chunk.Done {
				ch <- StreamChunk{Done: true}
				return
			}
		}
	}()

	return ch, nil
}

func (p *ollamaProvider) buildBody(req ProviderRequest, stream bool) map[string]any {
	msgs := make([]map[string]string, 0, len(req.Messages)+1)
	if req.System != "" {
		msgs = append(msgs, map[string]string{"role": "system", "content": req.System})
	}
	for _, m := range req.Messages {
		msgs = append(msgs, map[string]string{"role": m.Role, "content": m.Content})
	}
	return map[string]any{
		"model":    req.Model,
		"messages": msgs,
		"stream":   stream,
		"options": map[string]any{
			"temperature": req.Temperature,
			"num_predict": req.MaxTokens,
		},
	}
}

// --- helpers ---

func jsonBody(v any) io.Reader {
	b, _ := json.Marshal(v)
	return bytes.NewReader(b)
}
