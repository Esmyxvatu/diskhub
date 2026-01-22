package ollama

import (
	"bytes"
	"diskhub/web/config"
	"diskhub/web/logger"
	"diskhub/web/render"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strings"
)

var Messages []LLMMessage
var UIMessages []UIMessage

func ListModelAvaillable() []string {
	models := []string{}

	res, err := http.Get(fmt.Sprintf("http://localhost:%d/api/tags", config.Configuration.Ollama.Port))
	logger.Console.Verify(err)

	defer res.Body.Close()
	if res.StatusCode != 200 {
		logger.Console.Fatal("HTTP status %d wasn't expected", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	logger.Console.Verify(err)

	data := map[string][]Model{"models": {}}

	err = json.Unmarshal(body, &data)
	logger.Console.Verify(err)

	for key, value := range data {
		if key == "models" {
			for _, val := range value {
				models = append(models, val.Model)
			}
		}
	}

	logger.Console.Info("Models: %v", models)

	return models
}

func AskOllama(model string, content string, stream bool) UIMessage {
	Messages = append(Messages, LLMMessage{Role: "user", Content: content})
	UIMessages = append(UIMessages, UIMessage{Role: "user", Content: render.MarkdownToHTML([]byte(content))})
	question := Question{Model: model, Messages: Messages, Stream: stream}

	if stream {
		return UIMessage{Role: "system", Content: template.HTML("Stream response are'nt supported yet")}
	}

	body, err := json.Marshal(question)
	if err != nil {
		return UIMessage{Role: "system", Content: template.HTML(fmt.Sprintf("An error occured: %s", err.Error()))}
	}

	client := &http.Client{}

	req, err := http.NewRequest("POST", fmt.Sprintf("http://localhost:%d/api/chat", config.Configuration.Ollama.Port), bytes.NewBuffer(body))
	if err != nil {
		return UIMessage{Role: "system", Content: template.HTML(fmt.Sprintf("An error occured: %s", err.Error()))}
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return UIMessage{Role: "system", Content: template.HTML(fmt.Sprintf("An error occured: %s", err.Error()))}
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return UIMessage{Role: "system", Content: template.HTML(fmt.Sprintf("An error occured: %s", err.Error()))}
	}

	answer := Answer{}
	err = json.Unmarshal(resBody, &answer)
	if err != nil {
		return UIMessage{Role: "system", Content: template.HTML(fmt.Sprintf("An error occured: %s", err.Error()))}
	}
	if answer.Error != "" {
		return UIMessage{Role: "system", Content: template.HTML(answer.Error)}
	}

	raw := answer.Message.Content
	raw = template.HTMLEscapeString(raw)
	raw = strings.ReplaceAll(raw, "<think>", "```think\n")
	raw = strings.ReplaceAll(raw, "</think>", "\n```")

	Messages = append(Messages, answer.Message)
	UIMessages = append(UIMessages, UIMessage{
		Role:    "assistant",
		Content: render.MarkdownToHTML([]byte(raw)),
	})

	return UIMessages[len(UIMessages)-1]
}
