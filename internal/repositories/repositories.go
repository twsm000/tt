package repositories

import "github.com/twsm000/tt/internal/entities"

type TranslationRepository interface {
	FindByWord(word string) (*entities.Translation, error)
	Create(t *entities.Translation) error
	Update(t *entities.Translation) error
}
