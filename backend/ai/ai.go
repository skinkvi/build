package ai

import (
	"fmt"
	"os/exec"
)

type AIRequest struct {
	Model     string      `json:"model"`
	Messages  []AIMessage `json:"messages"`
	MaxTokens int         `json:"max_tokens"`
}

type AIMessage struct {
	Role    string      `json:"role"`
	Content interface{} `json:"content"`
}

type AIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func GetAIResponse(description string, imagePath string) (string, error) {
	cmd := exec.Command("python3", "./ai/ai.py", imagePath, description)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("Ошибка при выполнении Python-скрипта: %v\nВывод: %s", err, string(output))
	}
	return string(output), nil
}

func RunAI(imagePath string, promptText string) (string, error) {
	cmd := exec.Command("python3", "./ai/ai.py", imagePath, promptText)

	// Захватываем stdout и stderr
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("Ошибка при выполнении ai.py: %v\nВывод: %s", err, string(output))
	}
	return string(output), nil
}
