package author

import (
	"context"
	"fmt"
	"strings"

	"github.com/afteracademy/goserve/api/blog"
	"github.com/afteracademy/goserve/api/blog/dto"
	"github.com/afteracademy/goserve/api/blog/model"
	userModel "github.com/afteracademy/goserve/api/user/model"
	coredto "github.com/afteracademy/goserve/arch/dto"
	"github.com/afteracademy/goserve/arch/network"
	"github.com/afteracademy/goserve/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service interface {
	CreateBlog(createBlogDto *dto.BlogCreate, author *userModel.User) (*dto.BlogPrivate, error)
	UpdateBlog(updateBlogDto *dto.BlogUpdate, author *userModel.User) (*dto.BlogPrivate, error)
	DeactivateBlog(blogId uuid.UUID, author *userModel.User) error
	BlogSubmission(blogId uuid.UUID, author *userModel.User, submit bool) error
	GetBlogById(id uuid.UUID, author *userModel.User) (*dto.BlogPrivate, error)
	GetPaginatedDrafts(author *userModel.User, p *coredto.Pagination) ([]*dto.BlogInfo, error)
	GetPaginatedPublished(author *userModel.User, p *coredto.Pagination) ([]*dto.BlogInfo, error)
	GetPaginatedSubmitted(author *userModel.User, p *coredto.Pagination) ([]*dto.BlogInfo, error)
}

type service struct {
	network.BaseService
	db          *pgxpool.Pool
	blogService blog.Service
}

func NewService(db *pgxpool.Pool, blogService blog.Service) Service {
	return &service{
		BaseService: network.NewBaseService(),
		db:          db,
		blogService: blogService,
	}
}

func (s *service) CreateBlog(
	d *dto.BlogCreate,
	author *userModel.User,
) (*dto.BlogPrivate, error) {

	slug := utils.FormatEndpoint(d.Slug)
	exists := s.blogService.BlogSlugExists(slug)
	if exists {
		return nil, network.NewBadRequestError(
			"Blog with slug: "+slug+" already exists",
			nil,
		)
	}

	ctx := context.Background()
	var blog model.Blog

	query := `
		INSERT INTO blogs (
			title,
			description,
			draft_text,
			tags,
			author_id,
			img_url,
			slug
		)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7
		)
		RETURNING
			id,
			title,
			description,
			draft_text,
			tags,
			author_id,
			img_url,
			slug,
			score,
			submitted,
			drafted,
			published,
			status
	`

	err := s.db.QueryRow(
		ctx,
		query,
		d.Title,
		d.Description,
		d.DraftText,
		d.Tags,
		author.ID,
		d.ImgURL,
		slug,
	).Scan(
		&blog.ID,
		&blog.Title,
		&blog.Description,
		&blog.DraftText,
		&blog.Tags,
		&blog.AuthorID,
		&blog.ImgURL,
		&blog.Slug,
		&blog.Score,
		&blog.Submitted,
		&blog.Drafted,
		&blog.Published,
		&blog.Status,
	)

	if err != nil {
		return nil, err
	}

	return dto.NewBlogPrivate(&blog, author)
}

func (s *service) UpdateBlog(
	b *dto.BlogUpdate,
	author *userModel.User,
) (*dto.BlogPrivate, error) {
	ctx := context.Background()

	// Fetch existing blog (ownership + status check)
	selectQuery := `
		SELECT
			id,
			slug
		FROM blogs
		WHERE id = $1
		  AND author_id = $2
		  AND status = TRUE
	`

	var blogID uuid.UUID
	var currentSlug string

	err := s.db.QueryRow(
		ctx,
		selectQuery,
		b.ID,
		author.ID,
	).Scan(&blogID, &currentSlug)

	if err != nil {
		return nil, network.NewNotFoundError(
			"Blog with id: "+b.ID.String()+" does not exist",
			nil,
		)
	}

	// Build dynamic UPDATE
	setClauses := []string{}
	args := []any{}
	argPos := 1

	if b.Slug != nil {
		slug := utils.FormatEndpoint(*b.Slug)
		if slug != currentSlug {
			exists := s.blogService.BlogSlugExists(slug)
			if exists {
				return nil, network.NewBadRequestError(
					"Blog with slug: "+slug+" already exists",
					nil,
				)
			}
			setClauses = append(setClauses, fmt.Sprintf("slug = $%d", argPos))
			args = append(args, slug)
			argPos++
		}
	}

	if b.Title != nil {
		setClauses = append(setClauses, fmt.Sprintf("title = $%d", argPos))
		args = append(args, *b.Title)
		argPos++
	}

	if b.Description != nil {
		setClauses = append(setClauses, fmt.Sprintf("description = $%d", argPos))
		args = append(args, *b.Description)
		argPos++
	}

	if b.DraftText != nil {
		setClauses = append(setClauses, fmt.Sprintf("draft_text = $%d", argPos))
		args = append(args, *b.DraftText)
		argPos++
	}

	if b.Tags != nil {
		setClauses = append(setClauses, fmt.Sprintf("tags = $%d", argPos))
		args = append(args, *b.Tags)
		argPos++
	}

	if b.ImgURL != nil {
		setClauses = append(setClauses, fmt.Sprintf("img_url = $%d", argPos))
		args = append(args, *b.ImgURL)
		argPos++
	}

	// update timestamp
	setClauses = append(setClauses, "updated_at = CURRENT_TIMESTAMP")

	if len(setClauses) == 1 {
		// only updated_at then no meaningful change
		return s.GetBlogById(blogID, author)
	}

	updateQuery := fmt.Sprintf(`
		UPDATE blogs
		SET %s
		WHERE id = $%d
		  AND author_id = $%d
		  AND status = TRUE
	`,
		strings.Join(setClauses, ", "),
		argPos,
		argPos+1,
	)

	args = append(args, blogID, author.ID)

	tag, err := s.db.Exec(ctx, updateQuery, args...)
	if err != nil {
		return nil, err
	}

	if tag.RowsAffected() == 0 {
		return nil, network.NewNotFoundError("blog not found", nil)
	}

	// Return updated blog
	return s.GetBlogById(blogID, author)
}

func (s *service) DeactivateBlog(
	blogID uuid.UUID,
	author *userModel.User,
) error {
	ctx := context.Background()

	query := `
		UPDATE blogs
		SET
			status = FALSE,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		  AND author_id = $2
		  AND status = TRUE
	`

	tag, err := s.db.Exec(
		ctx,
		query,
		blogID,
		author.ID,
	)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return network.NewNotFoundError("blog not found", nil)
	}

	return nil
}

func (s *service) BlogSubmission(
	blogID uuid.UUID,
	author *userModel.User,
	submit bool,
) error {
	ctx := context.Background()

	query := `
		UPDATE blogs
		SET
			submitted = $1,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
		  AND author_id = $3
		  AND status = TRUE
	`

	tag, err := s.db.Exec(
		ctx,
		query,
		submit,
		blogID,
		author.ID,
	)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return network.NewNotFoundError("blog not found", nil)
	}

	return nil
}

func (s *service) GetBlogById(
	id uuid.UUID,
	author *userModel.User,
) (*dto.BlogPrivate, error) {
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
		  AND author_id = $2
		  AND status = TRUE
	`

	var b model.Blog

	err := s.db.QueryRow(ctx, query, id, author.ID).
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

	return dto.NewBlogPrivate(&b, author)
}

func (s *service) GetPaginatedDrafts(
	author *userModel.User,
	p *coredto.Pagination,
) ([]*dto.BlogInfo, error) {
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
		  AND drafted = TRUE
		  AND author_id = $1
		ORDER BY published_at DESC
		LIMIT $2 OFFSET $3
	`
	return s.getPaginated(query, author, p)
}

func (s *service) GetPaginatedPublished(
	author *userModel.User,
	p *coredto.Pagination,
) ([]*dto.BlogInfo, error) {
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
		  AND author_id = $1
		ORDER BY published_at DESC
		LIMIT $2 OFFSET $3
	`
	return s.getPaginated(query, author, p)
}

func (s *service) GetPaginatedSubmitted(
	author *userModel.User,
	p *coredto.Pagination,
) ([]*dto.BlogInfo, error) {
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
		  AND author_id = $1
		ORDER BY published_at DESC
		LIMIT $2 OFFSET $3
	`
	return s.getPaginated(query, author, p)
}

func (s *service) getPaginated(
	query string,
	author *userModel.User,
	p *coredto.Pagination,
) ([]*dto.BlogInfo, error) {

	ctx := context.Background()
	offset := (p.Page - 1) * p.Limit

	rows, err := s.db.Query(ctx, query, author.ID, p.Limit, offset)
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
