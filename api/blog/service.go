package blog

import (
	"context"
	"errors"
	"time"

	"github.com/afteracademy/goserve/api/blog/dto"
	"github.com/afteracademy/goserve/api/blog/model"
	"github.com/afteracademy/goserve/api/user"
	"github.com/afteracademy/goserve/arch/network"
	"github.com/afteracademy/goserve/arch/redis"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service interface {
	SetBlogDtoCacheById(blog *dto.BlogPublic) error
	GetBlogDtoCacheById(id uuid.UUID) (*dto.BlogPublic, error)
	SetBlogDtoCacheBySlug(blog *dto.BlogPublic) error
	GetBlogDtoCacheBySlug(slug string) (*dto.BlogPublic, error)
	BlogSlugExists(slug string) bool
	GetPublisedBlogById(id uuid.UUID) (*dto.BlogPublic, error)
	GetPublishedBlogBySlug(slug string) (*dto.BlogPublic, error)
}

type service struct {
	network.BaseService
	db              *pgxpool.Pool
	publicBlogCache redis.Cache[dto.BlogPublic]
	userService     user.Service
}

func NewService(db *pgxpool.Pool, store redis.Store, userService user.Service) Service {
	return &service{
		BaseService:     network.NewBaseService(),
		db:              db,
		publicBlogCache: redis.NewCache[dto.BlogPublic](store),
		userService:     userService,
	}
}

func (s *service) SetBlogDtoCacheById(blog *dto.BlogPublic) error {
	key := "blog_" + blog.ID.Hex()
	return s.publicBlogCache.SetJSON(key, blog, time.Duration(10*time.Minute))
}

func (s *service) GetBlogDtoCacheById(id uuid.UUID) (*dto.BlogPublic, error) {
	key := "blog_" + id.String()
	return s.publicBlogCache.GetJSON(key)
}

func (s *service) SetBlogDtoCacheBySlug(blog *dto.BlogPublic) error {
	key := "blog_" + blog.Slug
	return s.publicBlogCache.SetJSON(key, blog, time.Duration(10*time.Minute))
}

func (s *service) GetBlogDtoCacheBySlug(slug string) (*dto.BlogPublic, error) {
	key := "blog_" + slug
	return s.publicBlogCache.GetJSON(key)
}

func (s *service) BlogSlugExists(slug string) bool {

	query := `
		SELECT EXISTS (
			SELECT 1
			FROM blogs
			WHERE slug = $1
		)
	`

	var exists bool
	err := s.db.QueryRow(context.Background(), query, slug).Scan(&exists)
	if err != nil {
		return false
	}

	return exists
}

func (s *service) GetPublisedBlogById(blogID uuid.UUID) (*dto.BlogPublic, error) {
	ctx := context.Background()

	query := `
		SELECT
			id,
			title,
			description,
			text,
			slug,
			author_id,
			img_url,
			score,
			tags,
			published_at
		FROM blogs
		WHERE id = $1
		  AND status = TRUE
		  AND published = TRUE
	`

	var b model.Blog

	err := s.db.QueryRow(
		ctx,
		query,
		blogID,
	).Scan(
		&b.ID,
		&b.Title,
		&b.Description,
		&b.Text,
		&b.Slug,
		&b.AuthorID,
		&b.ImgURL,
		&b.Score,
		&b.Tags,
		&b.PublishedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, network.NewNotFoundError("blog not found", nil)
		}
		return nil, err
	}

	author, err := s.userService.FetchUserPublicProfile(b.AuthorID)
	if err != nil {
		return nil, network.NewNotFoundError("author not found", err)
	}

	return dto.NewBlogPublic(&b, author)
}

func (s *service) GetPublishedBlogBySlug(slug string) (*dto.BlogPublic, error) {
	ctx := context.Background()

	query := `
		SELECT
			id,
			title,
			description,
			text,
			slug,
			author_id,
			img_url,
			score,
			tags,
			published_at
		FROM blogs
		WHERE slug = $1
		  AND status = TRUE
		  AND published = TRUE
	`

	var b model.Blog

	err := s.db.QueryRow(
		ctx,
		query,
		slug,
	).Scan(
		&b.ID,
		&b.Title,
		&b.Description,
		&b.Text,
		&b.Slug,
		&b.AuthorID,
		&b.ImgURL,
		&b.Score,
		&b.Tags,
		&b.PublishedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, network.NewNotFoundError("blog not found", nil)
		}
		return nil, err
	}

	author, err := s.userService.FetchUserPublicProfile(b.AuthorID)
	if err != nil {
		return nil, network.NewNotFoundError("author not found", err)
	}

	return dto.NewBlogPublic(&b, author)
}
