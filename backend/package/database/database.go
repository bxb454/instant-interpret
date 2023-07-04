package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type Database struct {
	*sql.DB
}

type Message struct {
	ID                 int
	OriginalMessage    string
	OriginalLanguage   string
	TranslatedMessage  string
	TranslatedLanguage string
	Timestamp          time.Time
}

func NewDatabase(connStr string) (*Database, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("db.Ping: %v", err)
	}

	return &Database{db}, nil
}

func (db *Database) GetAllMessages() ([]Message, error) {
	rows, err := db.Query("SELECT id, original_message, original_language, translated_message, translated_language, timestamp FROM messages")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := []Message{}
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.ID, &msg.OriginalMessage, &msg.OriginalLanguage, &msg.TranslatedMessage, &msg.TranslatedLanguage, &msg.Timestamp); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, nil
}
