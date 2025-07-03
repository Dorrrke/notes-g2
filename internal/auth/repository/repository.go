package repository

import (
	"context"
	"errors"
	"time"

	"github.com/Dorrrke/notes-g2/internal/auth/models"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
)

type Repository struct {
	conn *pgx.Conn
}

func NewRepository(ctx context.Context, addr string) (*Repository, error) {
	conn, err := pgx.Connect(ctx, addr)
	if err != nil {
		return nil, err
	}
	return &Repository{conn: conn}, nil
}

func (db *Repository) SaveUser(user models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.conn.Exec(ctx, "INSERT INTO users(uid, name, email, password) VALUES ($1, $2, $3, $4)",
		user.UID, user.Name, user.Email, user.Password,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
				return ErrUserAlreadyExists
			}
		}
		return err
	}

	return nil
}

func (db *Repository) GetUser(email string) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	row := db.conn.QueryRow(ctx, "SELECT uid, name, email, password FROM users WHERE email = $1", email)
	var usr models.User

	if err := row.Scan(&usr.UID, &usr.Name, &usr.Email, &usr.Password); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, ErrUserNotFound
		}
		return models.User{}, err
	}

	return usr, nil
}
