// internal/imageprocessor/imageprocessor.go

package imageprocessor

import (
	"context"
	"os"
	"strings"

	"podbor/internal/config"

	vision "cloud.google.com/go/vision/apiv1"
	"go.uber.org/zap"
)

type ImageProcessor struct {
	logger *zap.Logger
	config *config.Config
}

func NewImageProcessor(logger *zap.Logger, config *config.Config) *ImageProcessor {
	return &ImageProcessor{
		logger: logger,
		config: config,
	}
}

func (ip *ImageProcessor) AnalyzeImage(imagePath string) (string, error) {
	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		ip.logger.Error("Ошибка при создании клиента Vision API", zap.Error(err))
		return "", err
	}
	defer client.Close()

	file, err := os.Open(imagePath)
	if err != nil {
		ip.logger.Error("Ошибка при открытии файла изображения", zap.Error(err))
		return "", err
	}
	defer file.Close()

	image, err := vision.NewImageFromReader(file)
	if err != nil {
		ip.logger.Error("Ошибка при создании изображения для Vision API", zap.Error(err))
		return "", err
	}

	annotations, err := client.DetectLabels(ctx, image, nil, 10)
	if err != nil {
		ip.logger.Error("Ошибка при анализе изображения", zap.Error(err))
		return "", err
	}

	var labels []string
	for _, annotation := range annotations {
		labels = append(labels, annotation.Description)
	}

	result := strings.Join(labels, ", ")
	return result, nil
}
