package migrations

import (
	"embed"
	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"io/fs"
)

//go:embed sql
var migrations embed.FS

func NewSQLProvider(conn *pgx.Conn) (*goose.Provider, error) {
	db, err := goose.OpenDBWithDriver("pgx", conn.Config().ConnString())
	if err != nil {
		return nil, err
	}

	fsys, err := fs.Sub(migrations, "sql")
	if err != nil {
		return nil, err
	}
	return goose.NewProvider(
		goose.DialectPostgres,
		db,
		fsys,
	)
}
