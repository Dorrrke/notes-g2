package dbstorage

import (
	"context"
	"time"

	notesDomain "github.com/Dorrrke/notes-g2/internal/domain/notes"
	"github.com/Dorrrke/notes-g2/pkg/logger"
)

func (db *DBStorage) GetNotes() ([]notesDomain.Note, error) {
	log := logger.Get()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := db.db.Query(ctx, "SELECT * FROM notes")
	if err != nil {
		log.Error().Err(err).Msg("failed to get notes")
		return nil, err
	}

	var notes []notesDomain.Note
	for rows.Next() {
		var note notesDomain.Note
		if err := rows.Scan(&note.NID, &note.Title, &note.Content, &note.Status, &note.Created_at, &note.UID); err != nil {
			log.Error().Err(err).Msg("failed to scan note")
			return nil, err
		}
		notes = append(notes, note)
	}

	return notes, nil
}

func (db *DBStorage) GetNote(nid string) (notesDomain.Note, error) {
	log := logger.Get()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var note notesDomain.Note

	row := db.db.QueryRow(ctx, "SELECT * FROM notes WHERE nid = $1", nid)
	err := row.Scan(&note.NID, &note.Title, &note.Content, &note.Status, &note.Created_at, &note.UID)
	if err != nil {
		log.Error().Err(err).Msg("failed to get note")
		return notesDomain.Note{}, err
	}

	return note, nil
}

func (db *DBStorage) SaveNotes(notes []notesDomain.Note) error {
	log := logger.Get()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := db.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			log.Error().Err(err).Msg("failed to rollback transaction")
		}
	}()

	_, err = tx.Prepare(
		ctx,
		"save_task",
		"INSERT INTO notes(nid, title, content, status, created_at, user_id) VALUES ($1, $2, $3, $4, $5, $6)",
	)
	if err != nil {
		return err
	}

	for _, note := range notes {
		_, err = tx.Exec(ctx, "save_task", note.NID, note.Title, note.Content, note.Status, note.Created_at, note.UID)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
