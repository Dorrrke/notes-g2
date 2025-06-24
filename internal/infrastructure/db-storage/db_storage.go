package dbstorage

import (
	"context"
	"errors"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5"
)

type DBStorage struct {
	db *pgx.Conn
}

func New(ctx context.Context, addr string) (*DBStorage, error) {
	conn, err := pgx.Connect(ctx, addr)
	if err != nil {
		return nil, err
	}

	return &DBStorage{db: conn}, nil
}

func (db *DBStorage) Close() error {
	return db.db.Close(context.Background())
}

func AppyMigrations(addr string) error {
	migrationPath := "file://migrations"
	m, err := migrate.New(migrationPath, addr)
	if err != nil {
		return err
	}

	defer m.Close()

	if err = m.Up(); err != nil && errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

// ! Применение миграций без использования пакета migrate
// func (db *DBStorage) AppyMigrations(ctx context.Context) error {
// 	currentVersion, err := db.checkDBVersion(ctx)
// 	if err != nil {
// 		return err
// 	}

// 	files, err := os.ReadDir("migrations")
// 	if err != nil {
// 		return err
// 	}

// 	for _, file := range files {
// 		if file.IsDir() || !strings.HasSuffix(file.Name(), ".up.sql") {
// 			continue
// 		}

// 		parts := strings.Split(file.Name(), "_")
// 		version, err := strconv.Atoi(parts[0])
// 		if err != nil {
// 			return err
// 		}

// 		if version > currentVersion {
// 			content, err := os.ReadFile(filepath.Join("migrations", file.Name()))
// 			if err != nil {
// 				return err
// 			}

// 			_, err = db.db.Exec(ctx, string(content))
// 			if err != nil {
// 				return err
// 			}

// 			_, err = db.db.Exec(ctx, `INSERT INTO schema_version (version) VALUES ($1)`, version)
// 			if err != nil {
// 				return err
// 			}
// 		}
// 	}

// 	return nil
// }

// func (db *DBStorage) checkDBVersion(ctx context.Context) (int, error) {
// 	log := logger.Get()
// 	var tableExist bool
// 	err := db.db.QueryRow(ctx, `SELECT EXISTS (
//             SELECT FROM information_schema.tables
//             WHERE table_name = 'schema_version'
//         )`).Scan(&tableExist)
// 	if err != nil {
// 		return -1, err
// 	}

// 	if !tableExist {
// 		log.Debug().Msg("schema_version table not found")
// 		return 0, nil
// 	}

// 	var version int
// 	err = db.db.QueryRow(ctx, `SELECT COALESCE(MAX(version), 0) FROM schema_version`).Scan(&version)
// 	if err != nil {
// 		return -1, err
// 	}

// 	return version, nil
// }
