package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"log"
)

type DataBase interface {
	Execute(query string, args ...interface{}) error
	Query(query string, args ...interface{}) (pgx.Rows, error)
	QueryRow(query string, args ...interface{}) pgx.Row
}

type Database struct {
	db *pgx.Conn
}

func NewDatabaseService(ctx context.Context) *Database {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:5432/%s", User, Password, Host, Name)
	db, err := pgx.Connect(ctx, connStr)
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v\n", err)
	}

	return &Database{db: db}
}

func (ds *Database) Close(ctx context.Context) {
	ds.db.Close(ctx)
}

func (ds *Database) Execute(ctx context.Context, query string, args ...interface{}) error {
	_, err := ds.db.Exec(ctx, query, args...)
	return err
}

func (ds *Database) Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	return ds.db.Query(ctx, query, args...)
}

func (ds *Database) QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	return ds.db.QueryRow(ctx, query, args...)
}
