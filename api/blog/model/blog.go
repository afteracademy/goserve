package model

import (
	"time"

	"github.com/google/uuid"
)

const BlogsTableName = "blogs"

type Blog struct {
	ID          uuid.UUID  // id
	Title       string     // title
	Description string     // description
	Text        *string    // text
	DraftText   string     // draft_text
	Tags        []string   // tags
	AuthorID    uuid.UUID  // author_id
	ImgURL      *string    // img_url
	Slug        string     // slug
	Score       float64    // score
	Views       int64      // views
	Likes       int64      // likes
	Comments    int64      // comments
	Flagged     bool       // flagged
	Submitted   bool       // submitted
	Drafted     bool       // drafted
	Published   bool       // published
	Status      bool       // status
	PublishedAt *time.Time // published_at
	CreatedAt   time.Time  // created_at
	UpdatedAt   time.Time  // updated_at
}
