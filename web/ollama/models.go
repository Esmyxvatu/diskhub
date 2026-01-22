package ollama

import "html/template"

type Model struct {
	Name       string         `json:"name"`
	Model      string         `json:"model"`
	LastMod    string         `json:"modified_at"`
	Size       int64          `json:"size"`
	Digest     string         `json:"digest"`
	Details    map[string]any `json:"details"`
	ParamSize  string         `json:"parameter_size"`
	QuantLevel string         `json:"quantization_level"`
}

type LLMMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type UIMessage struct {
	Role    string        `json:"role"`
	Content template.HTML `json:"content"`
}

type Question struct {
	Model    string       `json:"model"`
	Messages []LLMMessage `json:"messages"`
	Stream   bool         `json:"stream"`
}

type Answer struct {
	Error   string     `json:"error"`
	Model   string     `json:"model"`
	Message LLMMessage `json:"message"`
}

type UserQuestion struct {
	Stream  bool   `json:"stream"`
	Model   string `json:"model"`
	Content string `json:"content"`
}
