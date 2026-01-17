package model

import (
	"time"

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
