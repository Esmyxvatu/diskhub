package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"strings"
	"fmt"
)

type Model struct {
	Name		string			`json:"name"`
	Model		string			`json:"model"`
	LastMod		string			`json:"modified_at"`
	Size		int64			`json:"size"`
	Digest		string			`json:"digest"`
	Details		map[string]any		`json:"details"`
	ParamSize	string			`json:"parameter_size"`
	QuantLevel	string			`json:"quantization_level"`
}
type Message struct {
	Role	string		`json:"role"`
	Content template.HTML	`json:"content"`
}
type Question struct {
	Model	 string		`json:"model"`
	Messages []Message	`json:"messages"`
	Stream	 bool		`json:"stream"`
}
type Answer	struct {
	Model	string		`json:"model"`
	Message Message		`json:"message"`
}
type UserQuestion struct {
	Stream 	bool 	`json:"stream"`
	Model 	string 	`json:"model"`
	Content string 	`json:"content"`
}

var Messages []Message

func ListModelAvaillable() []string {
	models := []string{}

	res, err := http.Get( fmt.Sprintf("http://localhost:%d/api/tags", Configuration.Ollama.Port) )
	console.verify(err)

	defer res.Body.Close()
	if res.StatusCode != 200 { console.fatal("HTTP status %d wasn't expected", res.StatusCode) }

	body, err := io.ReadAll(res.Body)
	console.verify(err)

	data := map[string][]Model{ "models": {} }

	err = json.Unmarshal(body, &data)
	console.verify(err)

	for key, value := range data {
		if key == "models" {
			for _, val := range value {
				models = append(models, val.Model)
			}
		}
	}

	console.info("Models: %v", models)

	return models
}

func AskOllama(model string, content string, stream bool) Message {
	Messages = append(Messages, Message{ Role: "user", Content: template.HTML("<p>" + content + "</p>") })
	question := Question{ Model: model, Messages: Messages, Stream: stream}

	if stream {
		console.fatal("Stream response aren't currently supported")
	}

	body, err := json.Marshal(question)
	console.verify(err)

	client := &http.Client{}

	req, err := http.NewRequest("POST", fmt.Sprintf("http://localhost:%d/api/chat", Configuration.Ollama.Port), bytes.NewBuffer(body))
	console.verify(err)
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	console.verify(err)
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	console.verify(err)

	answer := Answer{}
	err = json.Unmarshal(resBody, &answer)
	console.verify(err)

	MessageContent := answer.Message.Content
	startIndex := strings.Index(string(MessageContent), "<think>")
	endIndex := strings.Index(string(MessageContent), "</think>")
	extractedThink := []byte(MessageContent[startIndex + 7 : endIndex])
	extractedContent := []byte(MessageContent[endIndex + 8 :])

	answer.Message.Content = template.HTML("<think>" + string(markdownToHTML(extractedThink)) + "</think>" + string(markdownToHTML(extractedContent)))

	Messages = append(Messages, answer.Message)

	return answer.Message
}
