package gen

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	openai "github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/williamnoble/kube-botany/pkg/plant"
	"io"
	"log/slog"
	"os"
	"strings"
)

type ImageGeneratorFunction func(plant string) error

type ImageGenerationService struct {
	staticDir string
	logger    *slog.Logger
	generator ImageGeneratorFunction
}

// NewImageGenerationService creates a new ImageGenerationService instance with the given staticDir and logger.
func NewImageGenerationService(staticDir string, logger *slog.Logger) *ImageGenerationService {
	s := ImageGenerationService{
		staticDir: staticDir,
		logger:    logger,
	}
	s.generator = s.GenerateImageOpenAI
	return &s
}

// NewMockImageGenerationService creates a new ImageGenerationService instance with the given staticDir and logger.
func NewMockImageGenerationService(
	staticDir string,
	logger *slog.Logger) *ImageGenerationService {
	s := ImageGenerationService{
		staticDir: staticDir,
		logger:    logger,
	}
	s.generator = s.GenerateMockImage
	s.logger.With("component", "generator").Info("configured mock image generation service")

	return &s
}

func (s *ImageGenerationService) ImageTask(plants map[string]*plant.Plant) error {
	var errs []error

	for _, p := range plants {
		plantImageName := p.Image()
		plantImagePath := fmt.Sprintf("%s/images/%s", s.staticDir, plantImageName)
		_, err := os.Stat(plantImagePath)

		if os.IsNotExist(err) {
			s.logger.With("component", "generator").Info("generating missing image", "image", plantImageName)
			err := s.generator(plantImageName)
			if err != nil {
				s.logger.With("component", "generator").Error("image generation failed", "image", plantImageName, "error", err.Error())
				errs = append(errs, fmt.Errorf("failed to generate image %s: %w", plantImageName, err))
				continue
			}
			s.logger.With("component", "generator").Info("image generated successfully", "image", plantImageName)
		} else if err != nil {
			s.logger.Error("failed to check if image exists", "image", plantImageName, "error", err.Error())
			errs = append(errs, fmt.Errorf("failed to check image %s: %w", plantImageName, err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("ImageTask encountered %d errors: %w", len(errs), errors.Join(errs...))
	}

	return nil
}

// GenerateImageOpenAI generates an image using OpenAI's ImageModelGPTImage1 model'.
func (s *ImageGenerationService) GenerateImageOpenAI(plant string) error {

	client := openai.NewClient(
		option.WithBaseURL("https://openrouter.ai/api/v1"))

	ctx := context.Background()
	prompt := "A robot ladybird in a forest of trees."

	image, err := client.Images.Generate(ctx, openai.ImageGenerateParams{
		Prompt:         prompt,
		Model:          openai.ImageModelGPTImage1,
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

	return nil
}

// GenerateMockImage uses a placeholder image to generate a mock image for a given plant.
func (s *ImageGenerationService) GenerateMockImage(plant string) error {
	plantName := strings.Split(plant, "-")[3]
	sourcePlaceholderImg := fmt.Sprintf("%s/%s", s.staticDir, fmt.Sprintf("0001-01-01-%s", plantName))
	destinationImg := fmt.Sprintf("%s/images/%s", s.staticDir, plant)
	src, err := os.Open(sourcePlaceholderImg)
	if err != nil {
		return fmt.Errorf("error opening source file: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(destinationImg)
	if err != nil {
		return fmt.Errorf("error opening destinationImg file: %w", err)
	}

	_, err = io.Copy(dst, src)
	if err != nil {
		return fmt.Errorf("error copying file contents: %w", err)
	}

	return nil
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
