package blogs

import (
	"context"
	"errors"
	"time"

	"github.com/afteracademy/goserve/api/blog/model"
	"github.com/afteracademy/goserve/api/blogs/dto"
	coredto "github.com/afteracademy/goserve/arch/dto"
	"github.com/afteracademy/goserve/arch/network"
	"github.com/afteracademy/goserve/arch/redis"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.mongodb.org/mongo-driver/bson"
)

type Service interface {
	SetSimilarBlogsDtoCache(blogId uuid.UUID, blogs []*dto.BlogItem) error
	GetSimilarBlogsDtoCache(blogId uuid.UUID) ([]*dto.BlogItem, error)
	GetPaginatedLatestBlogs(p *coredto.Pagination) ([]*dto.BlogItem, error)
	GetPaginatedTaggedBlogs(tag string, p *coredto.Pagination) ([]*dto.BlogItem, error)
	GetSimilarBlogs(blogId uuid.UUID) ([]*dto.BlogItem, error)
}

type service struct {
	network.BaseService
	db            *pgxpool.Pool
	itemBlogCache redis.Cache[dto.BlogItem]
}

func NewService(db *pgxpool.Pool, store redis.Store) Service {
	return &service{
		BaseService:   network.NewBaseService(),
		db:            db,
		itemBlogCache: redis.NewCache[dto.BlogItem](store),
	}
}

func (s *service) SetSimilarBlogsDtoCache(blogId uuid.UUID, blogs []*dto.BlogItem) error {
	key := "similar_blogs_" + blogId.String()
	return s.itemBlogCache.SetJSONList(key, blogs, 6*time.Hour)
}

func (s *service) GetSimilarBlogsDtoCache(blogId uuid.UUID) ([]*dto.BlogItem, error) {
	key := "similar_blogs_" + blogId.String()
	return s.itemBlogCache.GetJSONList(key)
}

func (s *service) GetPaginatedLatestBlogs(p *coredto.Pagination) ([]*dto.BlogItem, error) {
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
		ORDER BY published_at DESC, score DESC
		LIMIT $1 OFFSET $2
	`
	return s.getPaginated(query, p)
}

func (s *service) GetPaginatedTaggedBlogs(tag string, p *coredto.Pagination) ([]*dto.BlogItem, error) {
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
			AND $1 = ANY(tags)
		ORDER BY published_at DESC, score DESC
		LIMIT $2 OFFSET $3
	`
	ctx := context.Background()
	offset := (p.Page - 1) * p.Limit

	rows, err := s.db.Query(ctx, query, tag, p.Limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dtos []*dto.BlogItem

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

		d, err := dto.NewBlogItem(&b)
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

func (s *service) GetSimilarBlogs(
	blogID uuid.UUID,
) ([]*dto.BlogItem, error) {

	ctx := context.Background()
	var title string

	err := s.db.QueryRow(
		ctx,
		`
		SELECT title
		FROM blogs
		WHERE id = $1
		  AND published = TRUE
		  AND status = TRUE
		`,
		blogID,
	).Scan(&title)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, network.NewNotFoundError("blog not found", nil)
		}
		return nil, err
	}

	query := `
		SELECT
			id,
			title,
			description,
			slug,
			img_url,
			score,
			tags,
			published_at,
			ts_rank(
				to_tsvector('english', title),
				plainto_tsquery('english', $1)
			) AS similarity
		FROM blogs
		WHERE to_tsvector('english', title) @@ plainto_tsquery('english', $1)
		  AND published = TRUE
		  AND status = TRUE
		  AND id <> $2
		ORDER BY
			similarity DESC,
			updated_at DESC,
			score DESC
		LIMIT 6
	`

	rows, err := s.db.Query(ctx, query, title, blogID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*dto.BlogItem

	for rows.Next() {
		var (
			b          model.Blog
			similarity float32
		)

		if err := rows.Scan(
			&b.ID,
			&b.Title,
			&b.Description,
			&b.Slug,
			&b.ImgURL,
			&b.Score,
			&b.Tags,
			&b.PublishedAt,
			&similarity,
		); err != nil {
			return nil, err
		}

		item, err := dto.NewBlogItem(&b)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (s *service) getPublicPaginated(filter bson.M, p *coredto.Pagination) ([]*dto.BlogItem, error) {
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
		ORDER BY published_at DESC, score DESC
		LIMIT $1 OFFSET $2
	`
	return s.getPaginated(query, p)
}

func (s *service) getPaginated(
	query string,
	p *coredto.Pagination,
) ([]*dto.BlogItem, error) {

	ctx := context.Background()
	offset := (p.Page - 1) * p.Limit

	rows, err := s.db.Query(ctx, query, p.Limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dtos []*dto.BlogItem

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

		d, err := dto.NewBlogItem(&b)
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
