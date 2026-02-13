package repository

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/ram-ks/meeting-service/model"
)

type PreferredSlotRepository interface {
	Create(ctx context.Context, slot *model.PreferredSlot) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.PreferredSlot, error)
	GetByEmail(ctx context.Context, email string) ([]model.PreferredSlot, error)
	GetByEmails(ctx context.Context, emails []string) ([]model.PreferredSlot, error)
	Update(ctx context.Context, slot *model.PreferredSlot) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type preferredSlotRepository struct {
	db *sql.DB
}

func NewPreferredSlotRepository(db *sql.DB) PreferredSlotRepository {
	return &preferredSlotRepository{db: db}
}

func (r *preferredSlotRepository) Create(ctx context.Context, slot *model.PreferredSlot) error {
	query := `
		INSERT INTO preferred_slots (id, email, start_time, end_time, timezone, day_of_week, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.ExecContext(ctx, query,
		slot.ID, slot.Email,
		slot.StartTime, slot.EndTime, slot.Timezone, slot.DayOfWeek,
		slot.CreatedAt, slot.UpdatedAt,
	)
	return err
}

func (r *preferredSlotRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.PreferredSlot, error) {
	query := `
		SELECT id, email, start_time, end_time, timezone, day_of_week, created_at, updated_at
		FROM preferred_slots WHERE id = $1
	`
	slot := &model.PreferredSlot{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&slot.ID, &slot.Email,
		&slot.StartTime, &slot.EndTime, &slot.Timezone, &slot.DayOfWeek,
		&slot.CreatedAt, &slot.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return slot, nil
}

func (r *preferredSlotRepository) GetByEmail(ctx context.Context, email string) ([]model.PreferredSlot, error) {
	query := `
		SELECT id, email, start_time, end_time, timezone, day_of_week, created_at, updated_at
		FROM preferred_slots WHERE LOWER(email) = LOWER($1) ORDER BY start_time
	`
	rows, err := r.db.QueryContext(ctx, query, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var slots []model.PreferredSlot
	for rows.Next() {
		var slot model.PreferredSlot
		err := rows.Scan(
			&slot.ID, &slot.Email,
			&slot.StartTime, &slot.EndTime, &slot.Timezone, &slot.DayOfWeek,
			&slot.CreatedAt, &slot.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		slots = append(slots, slot)
	}
	return slots, nil
}

func (r *preferredSlotRepository) GetByEmails(ctx context.Context, emails []string) ([]model.PreferredSlot, error) {
	if len(emails) == 0 {
		return []model.PreferredSlot{}, nil
	}

	query := `
		SELECT id, email, start_time, end_time, timezone, day_of_week, created_at, updated_at
		FROM preferred_slots WHERE LOWER(email) = ANY($1) ORDER BY email, start_time
	`
	lowerEmails := make([]string, len(emails))
	for i, e := range emails {
		lowerEmails[i] = strings.ToLower(e)
	}

	rows, err := r.db.QueryContext(ctx, query, pq.Array(lowerEmails))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var slots []model.PreferredSlot
	for rows.Next() {
		var slot model.PreferredSlot
		err := rows.Scan(
			&slot.ID, &slot.Email,
			&slot.StartTime, &slot.EndTime, &slot.Timezone, &slot.DayOfWeek,
			&slot.CreatedAt, &slot.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		slots = append(slots, slot)
	}
	return slots, nil
}

func (r *preferredSlotRepository) Update(ctx context.Context, slot *model.PreferredSlot) error {
	query := `
		UPDATE preferred_slots 
		SET start_time = $1, end_time = $2, timezone = $3, day_of_week = $4, updated_at = $5 
		WHERE id = $6
	`
	slot.UpdatedAt = time.Now().UTC()
	_, err := r.db.ExecContext(ctx, query,
		slot.StartTime, slot.EndTime, slot.Timezone, slot.DayOfWeek, slot.UpdatedAt, slot.ID,
	)
	return err
}

func (r *preferredSlotRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM preferred_slots WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
