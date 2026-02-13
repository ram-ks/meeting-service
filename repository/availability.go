package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/ram-ks/meeting-service/model"
)

// Contract, Uppercase mean that it's public
type AvailabilityRepository interface {
	Create(ctx context.Context, availability *model.Availability) error
	Upsert(ctx context.Context, availability *model.Availability) error
	GetByEventID(ctx context.Context, eventID uuid.UUID) ([]model.Availability, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Availability, error)
	Update(ctx context.Context, availability *model.Availability) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// to implement an interface, one needs a type, this is it
type availabilityRepository struct {
	db *sql.DB
}

// Constructor function in go
func NewAvailabilityRepository(db *sql.DB) AvailabilityRepository {
	return &availabilityRepository{db: db}
}

// Mehod implementation
func (r *availabilityRepository) Create(ctx context.Context, availability *model.Availability) error {
	query := `
		INSERT INTO availability (id, event_id, participant_id, slot_id, status, available_from, available_to, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.ExecContext(ctx, query,
		availability.ID, availability.EventID, availability.ParticipantID, availability.SlotID,
		availability.Status, availability.AvailableFrom, availability.AvailableTo,
		availability.CreatedAt, availability.UpdatedAt,
	)
	return err
}

func (r *availabilityRepository) Upsert(ctx context.Context, availability *model.Availability) error {
	query := `
		INSERT INTO availability (id, event_id, participant_id, slot_id, status, available_from, available_to, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (participant_id, slot_id) DO UPDATE SET
			status = EXCLUDED.status,
			available_from = EXCLUDED.available_from,
			available_to = EXCLUDED.available_to,
			updated_at = EXCLUDED.updated_at
	`
	_, err := r.db.ExecContext(ctx, query,
		availability.ID, availability.EventID, availability.ParticipantID, availability.SlotID,
		availability.Status, availability.AvailableFrom, availability.AvailableTo,
		availability.CreatedAt, availability.UpdatedAt,
	)
	return err
}

func (r *availabilityRepository) GetByEventID(ctx context.Context, eventID uuid.UUID) ([]model.Availability, error) {
	query := `
		SELECT id, event_id, participant_id, slot_id, status, available_from, available_to, created_at, updated_at
		FROM availability WHERE event_id = $1
	`
	rows, err := r.db.QueryContext(ctx, query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var availabilities []model.Availability
	for rows.Next() {
		var a model.Availability
		err := rows.Scan(
			&a.ID, &a.EventID, &a.ParticipantID, &a.SlotID, &a.Status,
			&a.AvailableFrom, &a.AvailableTo, &a.CreatedAt, &a.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		availabilities = append(availabilities, a)
	}
	return availabilities, nil
}

func (r *availabilityRepository) Update(ctx context.Context, availability *model.Availability) error {
	query := `
		UPDATE availability SET status = $1, available_from = $2, available_to = $3, updated_at = $4
		WHERE id = $5
	`
	availability.UpdatedAt = time.Now().UTC()
	_, err := r.db.ExecContext(ctx, query,
		availability.Status, availability.AvailableFrom, availability.AvailableTo,
		availability.UpdatedAt, availability.ID,
	)
	return err
}

func (r *availabilityRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Availability, error) {
	query := `
		SELECT id, event_id, participant_id, slot_id, status, available_from, available_to, created_at, updated_at
		FROM availability WHERE id = $1
	`
	a := &model.Availability{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&a.ID, &a.EventID, &a.ParticipantID, &a.SlotID, &a.Status,
		&a.AvailableFrom, &a.AvailableTo, &a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (r *availabilityRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM availability WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
