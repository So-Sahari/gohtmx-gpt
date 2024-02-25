package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"Message": "GoHTMX GPT",
		})
	})
	router.POST("/run", func(c *gin.Context) {
		message := c.PostForm("prompt")
		promptMessage := fmt.Sprintf(`
You are a helpful AI assistant, who likes to make puns.
When asked to do something, you being with, 'you got it, boss!'.
The prompt is: %s
    `, message)
		response := callOllama(promptMessage)
		c.HTML(http.StatusOK, "prompt.html", gin.H{
			"Response": fmt.Sprintf("Assistant: %s", response),
		})
	})
	if err := router.Run(":31000"); err != nil {
		log.Fatal(err)
	}
}

func callOllama(prompt string) string {
	llm, err := ollama.New(ollama.WithModel("llama2"))
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	completion, err := llm.Call(ctx, fmt.Sprintf("Human: %s\nAssistant:", prompt),
		llms.WithTemperature(0.8),
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			fmt.Print(string(chunk))
			return nil
		}),
	)
	if err != nil {
		log.Fatal(err)
	}

	return completion
}
