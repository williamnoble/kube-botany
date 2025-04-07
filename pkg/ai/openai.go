package ai

import (
	"context"
	"encoding/base64"
	"github.com/openai/openai-go/option"
	"os"

	"github.com/openai/openai-go"
)

func GenerateImage() {
	client := openai.NewClient(
		option.WithBaseURL("https://openrouter.ai/api/v1"))

	ctx := context.Background()
	prompt := "A robot ladybird in a forest of trees."

	image, err := client.Images.Generate(ctx, openai.ImageGenerateParams{
		Prompt:         prompt,
		Model:          openai.ImageModelDallE3,
		ResponseFormat: openai.ImageGenerateParamsResponseFormatB64JSON,
		N:              openai.Int(1),
	})
	if err != nil {
		panic(err)
	}

	imageBytes, err := base64.StdEncoding.DecodeString(image.Data[0].B64JSON)
	if err != nil {
		panic(err)
	}

	dest := "./image.png"
	println("Writing image to " + dest)
	err = os.WriteFile(dest, imageBytes, 0755)
	if err != nil {
		panic(err)
	}
}

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
