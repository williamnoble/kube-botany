package server

import (
	"context"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

func GenerateText() {
	client := openai.NewClient(
		option.WithBaseURL("https://openrouter.ai/api/v1"))
	chatCompletion, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage("Say this is a test"),
		},
		Model: openai.ChatModelGPT4o,
	})
	if err != nil {
		panic(err.Error())
	}
	println(chatCompletion.Choices[0].Message.Content)
}
