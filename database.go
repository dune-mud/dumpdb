package main

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type database struct {
	db *sql.DB
}

var initialSchema = `
CREATE TABLE IF NOT EXISTS obj_dump (
  id           INTEGER PRIMARY KEY AUTOINCREMENT,
  name         TEXT NOT NULL,
  base_file    TEXT NOT NULL,
  size         INTEGER NOT NULL,
  full_size    INTEGER NOT NULL,
  refs         INTEGER NOT NULL DEFAULT 0,
  hb           INTEGER NOT NULL DEFAULT 0,
  environment  TEXT NOT NULL DEFAULT "",
  ticks        INTEGER NOT NULL DEFAULT 0,
  swap_status  TEXT NOT NULL DEFAULT "",
  created      TEXT NOT NULL
);
`

func openDB(path string) (*database, error) {
	db, err := sql.Open("sqlite3", path+"?cache=shared&mode=rwc&_sync=off&_journal=memory")
	if err != nil {
		return nil, fmt.Errorf("opening DB: %w", err)
	}

	_, err = db.Exec(initialSchema)
	if err != nil {
		return nil, fmt.Errorf("setting up DB schema: %w", err)
	}

	return &database{db: db}, nil
}

func (db database) Close() error {
	return db.db.Close()
}

func (db database) Insert(obs []*Object) error {
	tx, err := db.db.BeginTx(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}

	stmt, err := db.db.PrepareContext(context.Background(), `
  INSERT INTO obj_dump 
    (name, base_file, size, full_size, refs, hb, environment, ticks, swap_status, created) 
  VALUES 
    (   ?,         ?,    ?,         ?,    ?,  ?,           ?,     ?,           ?,       ?)`)
	if err != nil {
		return fmt.Errorf("preparing insert statement: %w", err)
	}

	defer stmt.Close()

	for _, o := range obs {
		o.Basefile = o.Name
		// Remove clone # from the `Name`, if present, to populate the `Basefile`
		if index := strings.Index(o.Basefile, "#"); index != -1 {
			o.Basefile = o.Basefile[0:index]
		}

		_, err = stmt.Exec(
			o.Name, o.Basefile, o.Size, o.FullSize, o.References, o.HB, o.Environment, o.Ticks, o.SwapStatus, o.Created)
		if err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("executing insert statement: %w", err)
		}
	}

	return tx.Commit()
}
