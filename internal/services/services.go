package services

import (
	"fmt"
	"os"
	"strings"

	"github.com/twsm000/tt/internal/entities"
	"github.com/twsm000/tt/internal/repositories"
)

type TranslationService interface {
	Translate(word string) (*entities.Translation, error)
}

func NewTranslationService(
	translator entities.Translator,
	repo repositories.TranslationRepository,
) TranslationService {
	return &translationService{
		translator: translator,
		repo:       repo,
	}
}

type translationService struct {
	translator entities.Translator
	repo       repositories.TranslationRepository
}

func (s *translationService) Translate(word string) (*entities.Translation, error) {
	word = strings.ToLower(word)
	t, err := s.repo.FindByWord(word)
	if err != nil {
		fmt.Fprintf(os.Stdout, "TranslationService: word %q not found: %v", word, err)
		translation, err := s.translator.Translate(word)
		if err != nil {
			return nil, err
		}
		var trans entities.Translation
		trans.Word = word
		trans.Translation = translation
		t = &trans
	}

	t.Count++
	switch t.ID {
	case 0:
		if err := s.repo.Create(t); err != nil {
			fmt.Fprintf(os.Stderr, "TranslationService: failed to create translation: %q - %v", word, err)
			return t, err
		}
	default:
		if err := s.repo.Update(t); err != nil {
			fmt.Fprintf(os.Stderr, "TranslationService: failed to update translation: %q - %v", word, err)
			return t, err
		}
	}
	return t, nil
}
