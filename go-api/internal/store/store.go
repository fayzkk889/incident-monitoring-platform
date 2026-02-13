package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LogEntry struct {
	ID        int64     `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Service   string    `json:"service"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	Metadata  string    `json:"metadata"`
}

type Incident struct {
	ID          int64      `json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	Status      string     `json:"status"`
	Severity    string     `json:"severity"`
	Description string     `json:"description"`
	Summary     *string    `json:"summary"`
	RootCause   *string    `json:"root_cause"`
	ResolvedAt  *time.Time `json:"resolved_at"`
}

type Repository interface {
	InsertLogs(ctx context.Context, logs []LogEntry) error
	ListRecentLogs(ctx context.Context, limit int) ([]LogEntry, error)

	CreateIncident(ctx context.Context, inc *Incident) error
	ListIncidents(ctx context.Context, limit int) ([]Incident, error)
	GetIncident(ctx context.Context, id int64) (*Incident, error)
	UpdateIncidentSummary(ctx context.Context, id int64, summary, rootCause string) error
}

type repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) Repository {
	return &repository{pool: pool}
}

func RunMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, `
CREATE TABLE IF NOT EXISTS logs (
    id SERIAL PRIMARY KEY,
    timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    service TEXT NOT NULL,
    level TEXT NOT NULL,
    message TEXT NOT NULL,
    metadata JSONB DEFAULT '{}'::jsonb
);

CREATE TABLE IF NOT EXISTS incidents (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    status TEXT NOT NULL DEFAULT 'open',
    severity TEXT NOT NULL DEFAULT 'medium',
    description TEXT NOT NULL,
    summary TEXT,
    root_cause TEXT,
    resolved_at TIMESTAMPTZ
);
`)
	return err
}

func (r *repository) InsertLogs(ctx context.Context, logs []LogEntry) error {
	batch := &pgx.Batch{}
	for _, l := range logs {
		batch.Queue(
			`INSERT INTO logs (timestamp, service, level, message, metadata)
             VALUES ($1, $2, $3, $4, COALESCE($5::jsonb, '{}'::jsonb))`,
			l.Timestamp, l.Service, l.Level, l.Message, l.Metadata,
		)
	}
	br := r.pool.SendBatch(ctx, batch)
	defer br.Close()
	_, err := br.Exec()
	return err
}

func (r *repository) ListRecentLogs(ctx context.Context, limit int) ([]LogEntry, error) {
	rows, err := r.pool.Query(ctx, `
SELECT id, timestamp, service, level, message, metadata
FROM logs
ORDER BY timestamp DESC
LIMIT $1
`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []LogEntry
	for rows.Next() {
		var l LogEntry
		if err := rows.Scan(&l.ID, &l.Timestamp, &l.Service, &l.Level, &l.Message, &l.Metadata); err != nil {
			return nil, err
		}
		res = append(res, l)
	}
	return res, rows.Err()
}

func (r *repository) CreateIncident(ctx context.Context, inc *Incident) error {
	return r.pool.QueryRow(ctx, `
INSERT INTO incidents (status, severity, description)
VALUES ($1, $2, $3)
RETURNING id, created_at
`, inc.Status, inc.Severity, inc.Description).Scan(&inc.ID, &inc.CreatedAt)
}

func (r *repository) ListIncidents(ctx context.Context, limit int) ([]Incident, error) {
	rows, err := r.pool.Query(ctx, `
SELECT id, created_at, status, severity, description, summary, root_cause, resolved_at
FROM incidents
ORDER BY created_at DESC
LIMIT $1
`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []Incident
	for rows.Next() {
		var inc Incident
		if err := rows.Scan(
			&inc.ID,
			&inc.CreatedAt,
			&inc.Status,
			&inc.Severity,
			&inc.Description,
			&inc.Summary,
			&inc.RootCause,
			&inc.ResolvedAt,
		); err != nil {
			return nil, err
		}
		res = append(res, inc)
	}
	return res, rows.Err()
}

func (r *repository) GetIncident(ctx context.Context, id int64) (*Incident, error) {
	row := r.pool.QueryRow(ctx, `
SELECT id, created_at, status, severity, description, summary, root_cause, resolved_at
FROM incidents
WHERE id = $1
`, id)

	var inc Incident
	if err := row.Scan(
		&inc.ID,
		&inc.CreatedAt,
		&inc.Status,
		&inc.Severity,
		&inc.Description,
		&inc.Summary,
		&inc.RootCause,
		&inc.ResolvedAt,
	); err != nil {
		return nil, err
	}
	return &inc, nil
}

func (r *repository) UpdateIncidentSummary(ctx context.Context, id int64, summary, rootCause string) error {
	_, err := r.pool.Exec(ctx, `
UPDATE incidents
SET summary = $2,
    root_cause = $3
WHERE id = $1
`, id, summary, rootCause)
	return err
}

