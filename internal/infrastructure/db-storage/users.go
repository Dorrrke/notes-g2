package dbstorage

import (
	"context"
	"errors"
	"time"

	usersDomain "github.com/Dorrrke/notes-g2/internal/domain/users"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

func (db *DBStorage) SaveUser(user usersDomain.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.db.Exec(ctx, "INSERT INTO users(id, name, email, password) VALUES ($1, $2, $3, $4)",
		user.UID, user.Name, user.Email, user.Password,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
				return usersDomain.ErrUserAlredyExists
			}
		}
		return err
	}

	return nil
}

func (db *DBStorage) GetUser(_ string) (usersDomain.User, error) {
	panic("unimplemented")
}
