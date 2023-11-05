package sqlite3

import (
	"database/sql"
	"errors"

	"github.com/twsm000/tt/internal/entities"
	"github.com/twsm000/tt/internal/repositories"
)

const (
	createTable = `CREATE TABLE IF NOT EXISTS translation (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		word TEXT,
		translation TEXT,
		count INTEGER
	);`
	findByWordQuery = `SELECT t.id, t.word, t.translation, t.count FROM translation t WHERE t.word = ?`
	insertIntoQuery = `INSERT INTO translation (word, translation, count) VALUES (?, ?, ?) RETURNING id`
	updateQuery     = `UPDATE translation SET word = ?, translation = ?, count = ? WHERE id = ?`
)

func NewTranslationRepository(db *sql.DB) (repositories.TranslationRepository, error) {
	if db == nil {
		return nil, errors.New("*sql.DB: invalid argument")
	}

	if _, err := db.Exec(createTable); err != nil {
		return nil, err
	}

	findByWordStmt, err := db.Prepare(findByWordQuery)
	if err != nil {
		return nil, err
	}
	insertIntoStmt, err := db.Prepare(insertIntoQuery)
	if err != nil {
		return nil, err
	}
	updateStmt, err := db.Prepare(updateQuery)
	if err != nil {
		return nil, err
	}

	return &translationRepository{
		db:             db,
		findByWordStmt: findByWordStmt,
		insertIntoStmt: insertIntoStmt,
		updateStmt:     updateStmt,
	}, nil
}

type translationRepository struct {
	db             *sql.DB
	findByWordStmt *sql.Stmt
	insertIntoStmt *sql.Stmt
	updateStmt     *sql.Stmt
}

func (r *translationRepository) FindByWord(word string) (*entities.Translation, error) {
	var t entities.Translation
	err := r.findByWordStmt.QueryRow(word).Scan(
		&t.ID,
		&t.Word,
		&t.Translation,
		&t.Count,
	)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *translationRepository) Create(t *entities.Translation) error {
	err := r.insertIntoStmt.QueryRow(t.Word, t.Translation, t.Count).Scan(&t.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *translationRepository) Update(t *entities.Translation) error {
	_, err := r.updateStmt.Exec(t.Word, t.Translation, t.Count, t.ID)
	if err != nil {
		return err
	}
	return nil
}
