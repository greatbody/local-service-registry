package store

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/greatbody/local-service-registry/internal/model"
	_ "github.com/mattn/go-sqlite3"
)

// Store manages service persistence via SQLite.
type Store struct {
	db *sql.DB
}

// New opens (or creates) the SQLite database at the given path and
// ensures the schema exists.
func New(dbPath string) (*Store, error) {
	db, err := sql.Open("sqlite3", dbPath+"?_journal_mode=WAL")
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}
	if err := migrate(db); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}
	return &Store{db: db}, nil
}

func migrate(db *sql.DB) error {
	const ddl = `
	CREATE TABLE IF NOT EXISTS services (
		id             TEXT PRIMARY KEY,
		name           TEXT NOT NULL,
		url            TEXT NOT NULL UNIQUE,
		description    TEXT NOT NULL DEFAULT '',
		remote_ip      TEXT NOT NULL DEFAULT '',
		status         TEXT NOT NULL DEFAULT 'unknown',
		registered_at  DATETIME NOT NULL,
		last_checked_at DATETIME
	);`
	if _, err := db.Exec(ddl); err != nil {
		return err
	}
	// Migration: add remote_ip column for existing databases.
	_, _ = db.Exec(`ALTER TABLE services ADD COLUMN remote_ip TEXT NOT NULL DEFAULT ''`)
	return nil
}

// Insert adds a new service record.
func (s *Store) Insert(svc *model.Service) error {
	const q = `INSERT INTO services (id, name, url, description, remote_ip, status, registered_at)
	            VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := s.db.Exec(q, svc.ID, svc.Name, svc.URL, svc.Description, svc.RemoteIP, svc.Status, svc.RegisteredAt)
	return err
}

// Delete removes a service by ID.
func (s *Store) Delete(id string) error {
	res, err := s.db.Exec(`DELETE FROM services WHERE id = ?`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("service %q not found", id)
	}
	return nil
}

// Get retrieves a single service by ID.
func (s *Store) Get(id string) (*model.Service, error) {
	const q = `SELECT id, name, url, description, remote_ip, status, registered_at, last_checked_at
	            FROM services WHERE id = ?`
	svc := &model.Service{}
	var lastChecked sql.NullTime
	err := s.db.QueryRow(q, id).Scan(
		&svc.ID, &svc.Name, &svc.URL, &svc.Description, &svc.RemoteIP,
		&svc.Status, &svc.RegisteredAt, &lastChecked,
	)
	if err != nil {
		return nil, err
	}
	if lastChecked.Valid {
		svc.LastCheckedAt = &lastChecked.Time
	}
	return svc, nil
}

// List returns all registered services.
func (s *Store) List() ([]*model.Service, error) {
	const q = `SELECT id, name, url, description, remote_ip, status, registered_at, last_checked_at
	            FROM services ORDER BY registered_at DESC`
	rows, err := s.db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []*model.Service
	for rows.Next() {
		svc := &model.Service{}
		var lastChecked sql.NullTime
		if err := rows.Scan(
			&svc.ID, &svc.Name, &svc.URL, &svc.Description, &svc.RemoteIP,
			&svc.Status, &svc.RegisteredAt, &lastChecked,
		); err != nil {
			return nil, err
		}
		if lastChecked.Valid {
			svc.LastCheckedAt = &lastChecked.Time
		}
		services = append(services, svc)
	}
	return services, rows.Err()
}

// UpdateStatus sets the health status and last-checked timestamp.
func (s *Store) UpdateStatus(id string, status model.HealthStatus, checkedAt time.Time) error {
	const q = `UPDATE services SET status = ?, last_checked_at = ? WHERE id = ?`
	_, err := s.db.Exec(q, status, checkedAt, id)
	return err
}

// Close closes the underlying database connection.
func (s *Store) Close() error {
	return s.db.Close()
}
