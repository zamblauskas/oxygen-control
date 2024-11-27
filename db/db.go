package db

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/zamblauskas/oxygen-control/models"
)

type DB struct {
	db *sql.DB
}

func NewDB(dbPath string) (*DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return &DB{db: db}, nil
}

func runMigrations(db *sql.DB) error {
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("could not create driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"sqlite3",
		driver,
	)
	if err != nil {
		return fmt.Errorf("could not create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("could not run migrations: %w", err)
	}

	return nil
}

func (d *DB) AddScheduleTrigger(t *models.ScheduleTrigger) error {
	id := uuid.New().String()

	query := `INSERT INTO triggers (id, type, schedule_hour, schedule_minute, schedule_action) 
						VALUES (?, ?, ?, ?, ?)`

	_, err := d.db.Exec(query, id, t.GetType(), t.Hour, t.Minute, t.Action)
	if err != nil {
		return err
	}

	return nil
}

func (d *DB) AddFlicTrigger(t *models.FlicTrigger) error {
	id := uuid.New().String()

	query := `INSERT INTO triggers (id, type, flic_button_mac, flic_button_on_single_click, flic_button_on_double_click, flic_button_on_hold) 
						VALUES (?, ?, ?, ?, ?, ?)`

	_, err := d.db.Exec(query, id, t.GetType(), t.Mac, t.OnSingleClick, t.OnDoubleClick, t.OnHold)
	if err != nil {
		return err
	}

	return nil
}
