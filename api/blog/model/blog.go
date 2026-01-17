package model

import (
	"time"

	"github.com/afteracademy/goserve/api/user/model"
	"github.com/google/uuid"
)

const BlogsTableName = "blogs"

type Blog struct {
	ID          uuid.UUID
	Title       string
	Description string
	Text        *string
	DraftText   string
	Tags        []string
	AuthorID    uuid.UUID
	ImgURL      *string
	Slug        string
	Score       float64
	Submitted   bool
	Drafted     bool
	Published   bool
	Status      bool
	PublishedAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewBlog(slug, title, description, draftText string, tags []string, author *model.User) *Blog {
	now := time.Now()
	b := Blog{
		Title:       title,
		Description: description,
		DraftText:   draftText,
		Tags:        tags,
		AuthorID:    author.ID,
		Slug:        slug,
		Score:       0.01,
		Submitted:   false,
		Drafted:     true,
		Published:   false,
		Status:      true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	return &b
}
