package translationservice

import (
	"context"
	"fmt"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
)

type TranslationService struct {
	client *translate.Client
}

func NewTranslationService(ctx context.Context) (*TranslationService, error) {
	client, err := translate.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("NewTranslationClient: %v", err)
	}

	return &TranslationService{
		client: client,
	}, nil
}

func (s *TranslationService) Close() error {
	return s.client.Close()
}

func (s *TranslationService) TranslateText(ctx context.Context, target language.Tag, text string) (string, error) {
	resp, err := s.client.Translate(ctx, []string{text}, target, nil)
	if err != nil {
		return "", fmt.Errorf("TranslateText: %v", err)
	}
	if len(resp) == 0 {
		return "", fmt.Errorf("TranslateText returned empty response to text: %s", text)
	}

	return resp[0].Text, nil
}

func (s *TranslationService) DetectLanguage(ctx context.Context, text string) (string, error) {
	lang, err := s.client.DetectLanguage(ctx, []string{text})
	if err != nil {
		return "", fmt.Errorf("DetectLanguage: %v", err)
	}
	if len(lang) == 0 || len(lang[0]) == 0 {
		return "", fmt.Errorf("DetectLanguage return value empty")
	}

	return lang[0][0].Language.String(), nil
}
