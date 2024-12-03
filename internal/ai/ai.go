package ai

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"podbor/internal/config"

	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
)

type AIService struct {
	logger *zap.Logger
	client *openai.Client
	config *config.Config
}

type AIResponse struct {
	Instructions string
	Materials    []string
	ImageURL     string // Если AI генерирует изображение
}

func NewAIService(logger *zap.Logger, config *config.Config) *AIService {
	client := openai.NewClient(config.AIService.APIKey)
	return &AIService{
		logger: logger,
		client: client,
		config: config,
	}
}

func (ai *AIService) ProcessIssue(description, imagePath string) (*AIResponse, error) {
	ctx := context.Background()

	// Чтение изображения и кодирование в base64
	imageData, err := os.ReadFile(imagePath)
	if err != nil {
		ai.logger.Error("Ошибка при чтении изображения", zap.Error(err))
		return nil, err
	}
	encodedImage := base64.StdEncoding.EncodeToString(imageData)
	imageBase64 := fmt.Sprintf("data:image/jpeg;base64,%s", encodedImage)

	// Составление промта
	userContent := fmt.Sprintf(`У меня возникла проблема: %s

На изображении ниже показано повреждение. Пожалуйста, предоставьте подробные пошаговые инструкции по починке этой проблемы и перечислите все необходимые материалы в виде списка. Верните ответ в формате JSON следующего вида:

{
  "instructions": "Ваши инструкции здесь",
  "materials": ["Материал1", "Материал2", "Материал3"]
}

Изображение: %s`, description, imageBase64)

	// Отправка запроса к ИИ
	req := openai.ChatCompletionRequest{
		Model: "gpt-3.5-turbo",
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    "user",
				Content: userContent,
			},
		},
		MaxTokens: 500,
	}

	// Отправка запроса к OpenAI API
	resp, err := ai.client.CreateChatCompletion(ctx, req)
	if err != nil {
		ai.logger.Error("Ошибка при вызове OpenAI API", zap.Error(err))
		return nil, err
	}

	if len(resp.Choices) == 0 {
		return nil, errors.New("нет ответов от OpenAI API")
	}

	// Парсинг ответа от AI
	content := resp.Choices[0].Message.Content
	instructions, materials, err := ai.parseAIResponse(content)
	if err != nil {
		ai.logger.Error("Ошибка при парсинге ответа AI", zap.Error(err))
		return nil, err
	}

	aiResponse := &AIResponse{
		Instructions: instructions,
		Materials:    materials,
	}

	return aiResponse, nil
}

func (ai *AIService) parseAIResponse(content string) (string, []string, error) {
	var response struct {
		Instructions string   `json:"instructions"`
		Materials    []string `json:"materials"`
	}

	err := json.Unmarshal([]byte(content), &response)
	if err != nil {
		ai.logger.Error("Ошибка при парсинге JSON от AI", zap.Error(err))
		return "", nil, err
	}

	return response.Instructions, response.Materials, nil
}
