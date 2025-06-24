package postgresql

// For now creating connections and operations in same file.
import (
	"database/sql"
	"log/slog"

	_ "github.com/lib/pq"
)

func NewPostgres(url string) (*sql.DB, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	slog.Info("message", "Database connection status : ", "successful")
	return db, nil
}
