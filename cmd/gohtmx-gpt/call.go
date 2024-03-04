package gohtmxgpt

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

var completion string

func callModel(c *gin.Context) {
	message := c.PostForm("prompt")
	promptMessage := fmt.Sprintf(`
You are a helpful AI assistant that only answers prompts if you know the answer.
The prompt is: %s
    `, message)

	input := completion + promptMessage
	response := callOllama(input)
	completion = input + response

	processString := strings.ReplaceAll(response, "\n", "<br>")
	processString = strings.ReplaceAll(processString, "    ", `<span class="tab-space"></span>`)
	c.HTML(http.StatusOK, "prompt.html", gin.H{
		"Response": template.HTML(processString),
	})
}

func callOllama(prompt string) string {
	llm, err := ollama.New(ollama.WithModel("llama2"))
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	response, err := llm.Call(ctx, fmt.Sprintf("Human: %s\nAssistant:", prompt),
		llms.WithTemperature(0.8),
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			fmt.Print(string(chunk))
			return nil
		}),
	)
	if err != nil {
		log.Fatal(err)
	}

	return response
}
