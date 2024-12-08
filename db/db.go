package db

import (
	"database/sql"
	"fmt"
	"time"

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

type TriggerDB struct {
	ID                      string
	Type                    string
	ScheduleHour            *int
	ScheduleMinute          *int
	ScheduleAction          *string
	FlicButtonMac           *string
	FlicButtonOnSingleClick *string
	FlicButtonOnDoubleClick *string
	FlicButtonOnHold        *string
	CreatedAt               time.Time
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

func (d *DB) GetAllTriggers() ([]models.Trigger, error) {
	query := `SELECT * FROM triggers`
	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var triggers []models.Trigger
	for rows.Next() {
		var t TriggerDB
		err := rows.Scan(
			&t.ID,
			&t.Type,
			&t.ScheduleHour,
			&t.ScheduleMinute,
			&t.ScheduleAction,
			&t.FlicButtonMac,
			&t.FlicButtonOnSingleClick,
			&t.FlicButtonOnDoubleClick,
			&t.FlicButtonOnHold,
			&t.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		// Convert DB model to appropriate trigger type
		if t.Type == string(models.TriggerTypeSchedule) {
			triggers = append(triggers, &models.ScheduleTrigger{
				ID:     t.ID,
				Type:   models.TriggerTypeSchedule,
				Hour:   *t.ScheduleHour,
				Minute: *t.ScheduleMinute,
				Action: models.TriggerAction(*t.ScheduleAction),
			})
		} else if t.Type == string(models.TriggerTypeFlic) {
			triggers = append(triggers, &models.FlicTrigger{
				ID:            t.ID,
				Type:          models.TriggerTypeFlic,
				Mac:           *t.FlicButtonMac,
				OnSingleClick: (*models.TriggerAction)(t.FlicButtonOnSingleClick),
				OnDoubleClick: (*models.TriggerAction)(t.FlicButtonOnDoubleClick),
				OnHold:        (*models.TriggerAction)(t.FlicButtonOnHold),
			})
		}
	}
	return triggers, nil
}
