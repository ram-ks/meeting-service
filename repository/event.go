package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/ram-ks/meeting-service/model"
)

type EventRepository interface {
	Create(ctx context.Context, event *model.Event) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Event, error)
	List(ctx context.Context, organizerID uuid.UUID) ([]model.Event, error)
	Update(ctx context.Context, event *model.Event) error
	Delete(ctx context.Context, id uuid.UUID) error
	CreateSlot(ctx context.Context, slot *model.TimeSlot) error
	GetSlotsByEventID(ctx context.Context, eventID uuid.UUID) ([]model.TimeSlot, error)
	GetSlotByID(ctx context.Context, id uuid.UUID) (*model.TimeSlot, error)
	UpdateSlot(ctx context.Context, slot *model.TimeSlot) error
	DeleteSlot(ctx context.Context, id uuid.UUID) error
	CreateParticipant(ctx context.Context, participant *model.Participant) error
	GetParticipantsByEventID(ctx context.Context, eventID uuid.UUID) ([]model.Participant, error)
	GetParticipantByID(ctx context.Context, id uuid.UUID) (*model.Participant, error)
	UpdateParticipantStatus(ctx context.Context, id uuid.UUID, status model.ParticipantStatus) error
}

type eventRepository struct {
	db *sql.DB
}

func NewEventRepository(db *sql.DB) EventRepository {
	return &eventRepository{db: db}
}

func (r *eventRepository) Create(ctx context.Context, event *model.Event) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO events (id, title, description, organizer_id, duration, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err = tx.ExecContext(ctx, query,
		event.ID, event.Title, event.Description, event.OrganizerID,
		event.Duration, event.Status, event.CreatedAt, event.UpdatedAt,
	)
	if err != nil {
		return err
	}

	for _, slot := range event.ProposedSlots {
		slotQuery := `
			INSERT INTO time_slots (id, event_id, start_time, end_time, timezone, created_at)
			VALUES ($1, $2, $3, $4, $5, $6)
		`
		_, err = tx.ExecContext(ctx, slotQuery,
			slot.ID, event.ID, slot.StartTime, slot.EndTime, slot.Timezone, slot.CreatedAt,
		)
		if err != nil {
			return err
		}
	}

	for _, participant := range event.Participants {
		participantQuery := `
			INSERT INTO participants (id, event_id, email, name, status, created_at)
			VALUES ($1, $2, $3, $4, $5, $6)
		`
		_, err = tx.ExecContext(ctx, participantQuery,
			participant.ID, event.ID, participant.Email,
			participant.Name, participant.Status, participant.CreatedAt,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *eventRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Event, error) {
	query := `
		SELECT id, title, description, organizer_id, duration, status, finalized_slot_id, created_at, updated_at
		FROM events WHERE id = $1
	`
	event := &model.Event{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&event.ID, &event.Title, &event.Description, &event.OrganizerID,
		&event.Duration, &event.Status, &event.FinalizedSlotID, &event.CreatedAt, &event.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	slots, err := r.GetSlotsByEventID(ctx, id)
	if err != nil {
		return nil, err
	}
	event.ProposedSlots = slots

	participants, err := r.GetParticipantsByEventID(ctx, id)
	if err != nil {
		return nil, err
	}
	event.Participants = participants

	return event, nil
}

func (r *eventRepository) List(ctx context.Context, organizerID uuid.UUID) ([]model.Event, error) {
	query := `
		SELECT id, title, description, organizer_id, duration, status, finalized_slot_id, created_at, updated_at
		FROM events WHERE organizer_id = $1 ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, organizerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []model.Event
	for rows.Next() {
		var event model.Event
		err := rows.Scan(
			&event.ID, &event.Title, &event.Description, &event.OrganizerID,
			&event.Duration, &event.Status, &event.FinalizedSlotID, &event.CreatedAt, &event.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}

func (r *eventRepository) Update(ctx context.Context, event *model.Event) error {
	query := `
		UPDATE events SET title = $1, description = $2, duration = $3, status = $4, 
		finalized_slot_id = $5, updated_at = $6 WHERE id = $7
	`
	event.UpdatedAt = time.Now().UTC()
	_, err := r.db.ExecContext(ctx, query,
		event.Title, event.Description, event.Duration, event.Status,
		event.FinalizedSlotID, event.UpdatedAt, event.ID,
	)
	return err
}

func (r *eventRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM events WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *eventRepository) CreateSlot(ctx context.Context, slot *model.TimeSlot) error {
	query := `
		INSERT INTO time_slots (id, event_id, start_time, end_time, timezone, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(ctx, query,
		slot.ID, slot.EventID, slot.StartTime, slot.EndTime, slot.Timezone, slot.CreatedAt,
	)
	return err
}

func (r *eventRepository) GetSlotsByEventID(ctx context.Context, eventID uuid.UUID) ([]model.TimeSlot, error) {
	query := `
		SELECT id, event_id, start_time, end_time, timezone, created_at
		FROM time_slots WHERE event_id = $1 ORDER BY start_time
	`
	rows, err := r.db.QueryContext(ctx, query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var slots []model.TimeSlot
	for rows.Next() {
		var slot model.TimeSlot
		err := rows.Scan(&slot.ID, &slot.EventID, &slot.StartTime, &slot.EndTime, &slot.Timezone, &slot.CreatedAt)
		if err != nil {
			return nil, err
		}
		slots = append(slots, slot)
	}
	return slots, nil
}

func (r *eventRepository) GetSlotByID(ctx context.Context, id uuid.UUID) (*model.TimeSlot, error) {
	query := `
		SELECT id, event_id, start_time, end_time, timezone, created_at
		FROM time_slots WHERE id = $1
	`
	slot := &model.TimeSlot{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&slot.ID, &slot.EventID, &slot.StartTime, &slot.EndTime, &slot.Timezone, &slot.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return slot, nil
}

func (r *eventRepository) UpdateSlot(ctx context.Context, slot *model.TimeSlot) error {
	query := `
		UPDATE time_slots SET start_time = $1, end_time = $2, timezone = $3 WHERE id = $4
	`
	_, err := r.db.ExecContext(ctx, query, slot.StartTime, slot.EndTime, slot.Timezone, slot.ID)
	return err
}

func (r *eventRepository) DeleteSlot(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM time_slots WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *eventRepository) CreateParticipant(ctx context.Context, participant *model.Participant) error {
	query := `
		INSERT INTO participants (id, event_id, email, name, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(ctx, query,
		participant.ID, participant.EventID,
		participant.Email, participant.Name, participant.Status, participant.CreatedAt,
	)
	return err
}

func (r *eventRepository) GetParticipantsByEventID(ctx context.Context, eventID uuid.UUID) ([]model.Participant, error) {
	query := `
		SELECT id, event_id, email, name, status, created_at
		FROM participants WHERE event_id = $1
	`
	rows, err := r.db.QueryContext(ctx, query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var participants []model.Participant
	for rows.Next() {
		var p model.Participant
		err := rows.Scan(&p.ID, &p.EventID, &p.Email, &p.Name, &p.Status, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		participants = append(participants, p)
	}
	return participants, nil
}

func (r *eventRepository) GetParticipantByID(ctx context.Context, id uuid.UUID) (*model.Participant, error) {
	query := `
		SELECT id, event_id, email, name, status, created_at
		FROM participants WHERE id = $1
	`
	p := &model.Participant{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.EventID, &p.Email, &p.Name, &p.Status, &p.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *eventRepository) UpdateParticipantStatus(ctx context.Context, id uuid.UUID, status model.ParticipantStatus) error {
	query := `UPDATE participants SET status = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, status, id)
	return err
}
