package contact

import (
	"context"

	"github.com/afteracademy/goserve/api/contact/dto"
	"github.com/afteracademy/goserve/api/contact/model"
	coredto "github.com/afteracademy/goserve/arch/dto"
	"github.com/afteracademy/goserve/arch/network"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service interface {
	SaveMessage(d *dto.MessageCreate) (*model.Message, error)
	FetchMessage(id uuid.UUID) (*model.Message, error)
	FetchPaginatedMessage(p *coredto.Pagination) ([]*model.Message, error)
}

type service struct {
	network.BaseService
	db *pgxpool.Pool
}

func NewService(db *pgxpool.Pool) Service {
	return &service{
		BaseService: network.NewBaseService(),
		db:          db,
	}
}

func (s *service) SaveMessage(d *dto.MessageCreate) (*model.Message, error) {
	return s.CreateMessage(context.Background(), d)
}

func (s *service) FetchMessage(id uuid.UUID) (*model.Message, error) {
	return s.FindMessage(context.Background(), id)
}

func (s *service) FetchPaginatedMessage(p *coredto.Pagination) ([]*model.Message, error) {
	return s.FindPaginatedMessage(context.Background(), p)
}

func (s *service) CreateMessage(
	ctx context.Context,
	dto *dto.MessageCreate,
) (*model.Message, error) {

	msg := model.Message{}

	query := `
		INSERT INTO messages (
			type,
			msg,
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING
			id,
			type,
			msg,
			status,
			created_at,
			updated_at
	`

	err := s.db.QueryRow(
		ctx,
		query,
		dto.Type,
		dto.Msg,
	).Scan(
		&msg.ID,
		&msg.Type,
		&msg.Msg,
		&msg.Status,
		&msg.CreatedAt,
		&msg.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &msg, nil
}

func (s *service) FindMessage(
	ctx context.Context,
	id uuid.UUID,
) (*model.Message, error) {

	query := `
		SELECT
			id,
			type,
			msg,
			status,
			created_at,
			updated_at
		FROM messages
		WHERE id = $1
	`

	var m model.Message

	err := s.db.QueryRow(ctx, query, id).
		Scan(
			&m.ID,
			&m.Type,
			&m.Msg,
			&m.Status,
			&m.CreatedAt,
			&m.UpdatedAt,
		)

	if err != nil {
		return nil, err
	}

	return &m, nil
}

func (s *service) FindPaginatedMessage(
	ctx context.Context,
	p *coredto.Pagination,
) ([]*model.Message, error) {

	offset := (p.Page - 1) * p.Limit

	query := `
		SELECT
			id,
			type,
			msg,
			status,
			created_at,
			updated_at
		FROM messages
		WHERE status = TRUE
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := s.db.Query(ctx, query, p.Limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*model.Message

	for rows.Next() {
		var m model.Message
		if err := rows.Scan(
			&m.ID,
			&m.Type,
			&m.Msg,
			&m.Status,
			&m.CreatedAt,
			&m.UpdatedAt,
		); err != nil {
			return nil, err
		}
		messages = append(messages, &m)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}
