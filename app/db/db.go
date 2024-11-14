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

func NewDatabaseService() *Database {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:5432/%s", User, Password, Host, Name)
	db, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v\n", err)
	}

	return &Database{db: db}
}

func (ds *Database) Close() {
	ds.db.Close(context.Background())
}

func (ds *Database) Execute(query string, args ...interface{}) error {
	_, err := ds.db.Exec(context.Background(), query, args...)
	return err
}

func (ds *Database) Query(query string, args ...interface{}) (pgx.Rows, error) {
	return ds.db.Query(context.Background(), query, args...)
}

func (ds *Database) QueryRow(query string, args ...interface{}) pgx.Row {
	return ds.db.QueryRow(context.Background(), query, args...)
}
