package main

import (
	sql "database/sql"
	"log"

	"sync"

	_ "github.com/mattn/go-sqlite3"
)

type Repository interface {
	Migrate() (sql.Result, error)
	Add(rec *Item) (sql.Result, error)
}

type SQLiteRepository struct {
	mut sync.Mutex
	db  *sql.DB
}

func New(database string) *SQLiteRepository {
	db, err := sql.Open("sqlite3", database)
	if err != nil {
		log.Fatal(err)
	}
	return &SQLiteRepository{
		db: db,
	}
}

func (r *SQLiteRepository) Migrate() (sql.Result, error) {
	query := `
		DROP TABLE IF EXISTS posts;
		CREATE TABLE IF NOT EXISTS posts (
				id INTEGER PRIMARY KEY ,
				author TEXT NOT NULL,
				type TEXT NOT NULL,
				text TEXT NOT NULL
		);
	`
	r.mut.Lock()
	defer r.mut.Unlock()
	return r.db.Exec(query)
}

func (r *SQLiteRepository) Add(rec Item) (sql.Result, error) {
	query := `
		INSERT INTO posts(id, author, type, text)
		VALUES(?, ?, ?, ?);
	`
	r.mut.Lock()
	defer r.mut.Unlock()
	return r.db.Exec(query, rec.Id, rec.By, rec.Type, rec.Text)
}
