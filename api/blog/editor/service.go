package editor

import (
	"context"
	"errors"
	"time"

	"github.com/afteracademy/goserve/api/blog/dto"
	"github.com/afteracademy/goserve/api/blog/model"
	"github.com/afteracademy/goserve/api/user"
	coredto "github.com/afteracademy/goserve/arch/dto"
	"github.com/afteracademy/goserve/arch/network"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service interface {
	GetBlogById(id uuid.UUID) (*dto.BlogPrivate, error)
	BlogPublication(blogId uuid.UUID, publish bool) error
	GetPaginatedPublished(p *coredto.Pagination) ([]*dto.BlogInfo, error)
	GetPaginatedSubmitted(p *coredto.Pagination) ([]*dto.BlogInfo, error)
}

type service struct {
	network.BaseService
	db          *pgxpool.Pool
	userService user.Service
}

func NewService(db *pgxpool.Pool, userService user.Service) Service {
	return &service{
		BaseService: network.NewBaseService(),
		db:          db,
		userService: userService,
	}
}

func (s *service) BlogPublication(
	blogID uuid.UUID,
	publish bool,
) error {
	ctx := context.Background()

	selectQuery := `
		SELECT
			published,
			submitted,
			drafted,
			draft_text,
			published_at
		FROM blogs
		WHERE id = $1
		  AND status = TRUE
	`

	var (
		published   bool
		submitted   bool
		drafted     bool
		draftText   string
		publishedAt *time.Time
	)

	err := s.db.QueryRow(
		ctx,
		selectQuery,
		blogID,
	).Scan(
		&published,
		&submitted,
		&drafted,
		&draftText,
		&publishedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return network.NewNotFoundError(
				"blog for id "+blogID.String()+" not found",
				nil,
			)
		}
		return err
	}

	if publish {
		if published {
			return network.NewBadRequestError(
				"blog for id "+blogID.String()+" is already published",
				nil,
			)
		}
		if !submitted {
			return network.NewBadRequestError(
				"blog for id "+blogID.String()+" is not submitted",
				nil,
			)
		}
	} else {
		if !published {
			return network.NewBadRequestError(
				"blog for id "+blogID.String()+" is not published",
				nil,
			)
		}
	}

	updateQuery := `
		UPDATE blogs
		SET
			drafted = $1,
			submitted = $2,
			published = $3,
			text = $4,
			published_at = $5,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $6
		  AND status = TRUE
	`

	var text *string

	if publish {
		if publishedAt == nil {
			now := time.Now()
			publishedAt = &now
		}
		text = &draftText
	} else {
		publishedAt = nil
		text = nil
	}

	tag, err := s.db.Exec(
		ctx,
		updateQuery,
		!publish,
		false,
		publish,
		text,
		publishedAt,
		blogID,
	)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return network.NewNotFoundError("blog not found", nil)
	}

	return nil
}

func (s *service) GetBlogById(id uuid.UUID) (*dto.BlogPrivate, error) {
	ctx := context.Background()

	query := `
		SELECT
			id,
			title,
			description,
			text,
			draft_text,
			tags,
			author_id,
			img_url,
			slug,
			score,
			submitted,
			drafted,
			published,
			status,
			published_at,
			created_at,
			updated_at
		FROM blogs
		WHERE id = $1
		  AND status = TRUE
	`

	var b model.Blog

	err := s.db.QueryRow(ctx, query, id).
		Scan(
			&b.ID,
			&b.Title,
			&b.Description,
			&b.Text,
			&b.DraftText,
			&b.Tags,
			&b.AuthorID,
			&b.ImgURL,
			&b.Slug,
			&b.Score,
			&b.Submitted,
			&b.Drafted,
			&b.Published,
			&b.Status,
			&b.PublishedAt,
			&b.CreatedAt,
			&b.UpdatedAt,
		)

	if err != nil {
		return nil, err
	}

	author, err := s.userService.FetchUserById(b.AuthorID)
	if err != nil {
		return nil, network.NewInternalServerError("failed to fetch author", err)
	}
	if author == nil {
		return nil, network.NewNotFoundError("author not found", nil)
	}

	return dto.NewBlogPrivate(&b, author)
}

func (s *service) GetPaginatedPublished(p *coredto.Pagination) ([]*dto.BlogInfo, error) {
	query := `
		SELECT
			id,
			title,
			description,
			slug,
			img_url,
			score,
			tags,
			published_at
		FROM blogs
		WHERE status = TRUE
		  AND published = TRUE
		ORDER BY published_at DESC
		LIMIT $1 OFFSET $2
	`
	return s.getPaginated(query, p)
}

func (s *service) GetPaginatedSubmitted(p *coredto.Pagination) ([]*dto.BlogInfo, error) {
	query := `
		SELECT
			id,
			title,
			description,
			slug,
			img_url,
			score,
			tags,
			published_at
		FROM blogs
		WHERE status = TRUE
		  AND submitted = TRUE
		ORDER BY published_at DESC
		LIMIT $1 OFFSET $2
	`
	return s.getPaginated(query, p)
}

func (s *service) getPaginated(
	query string,
	p *coredto.Pagination,
) ([]*dto.BlogInfo, error) {

	ctx := context.Background()
	offset := (p.Page - 1) * p.Limit

	rows, err := s.db.Query(ctx, query, p.Limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dtos []*dto.BlogInfo

	for rows.Next() {
		var b model.Blog
		if err := rows.Scan(
			&b.ID,
			&b.Title,
			&b.Description,
			&b.Slug,
			&b.ImgURL,
			&b.Score,
			&b.Tags,
			&b.PublishedAt,
		); err != nil {
			return nil, err
		}

		d, err := dto.NewBlogInfo(&b)
		if err != nil {
			return nil, err
		}

		dtos = append(dtos, d)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return dtos, nil
}
