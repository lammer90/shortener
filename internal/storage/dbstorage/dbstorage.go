package dbstorage

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/lammer90/shortener/internal/models"
	"github.com/lammer90/shortener/internal/storage"
)

type dbStorage struct {
	db *sql.DB
}

func New(db *sql.DB) dbStorage {
	initDB(db)
	return dbStorage{db: db}
}

func (d dbStorage) Save(key, value string) error {
	_, err := d.db.ExecContext(context.Background(), `
        INSERT INTO shorts
        (short_url, original_url)
        VALUES
        ($1, $2);
    `, key, value)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			shortURL := findByOriginal(d.db, value)
			return storage.NewErrConflict(shortURL, err)
		}
	}
	return err
}

func (d dbStorage) SaveBatch(shorts []*models.BatchToSave) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	ctx := context.Background()
	stmt, err := d.db.PrepareContext(ctx, `
        INSERT INTO shorts
        (short_url, original_url)
        VALUES
        ($1, $2)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, short := range shorts {
		_, err := stmt.ExecContext(ctx, short.ShortURL, short.OriginalURL)
		if err != nil {
			return err
		}
	}

	tx.Commit()
	return nil
}

func (d dbStorage) Find(key string) (string, bool, error) {
	row := d.db.QueryRowContext(context.Background(), `
        SELECT
            s.original_url
        FROM shorts s
        WHERE
            s.short_url = $1
    `, key)

	var value string
	err := row.Scan(&value)
	if err != nil {
		return "", false, err
	}
	return value, true, nil
}

func initDB(db *sql.DB) {
	ctx := context.Background()
	db.ExecContext(ctx, `
        CREATE TABLE IF NOT EXISTS shorts (
            short_url varchar,
            original_url varchar
        )
    `)
	db.ExecContext(ctx, `
        CREATE UNIQUE INDEX IF NOT EXISTS shorts_original_url_idx ON shorts (original_url)
    `)
}

func findByOriginal(db *sql.DB, value string) string {
	row := db.QueryRowContext(context.Background(), `
        SELECT
            s.short_url
        FROM shorts s
        WHERE
            s.original_url = $1
    `, value)

	var shortURL string
	err := row.Scan(&shortURL)
	if err != nil {
		return ""
	}
	return shortURL
}
