CREATE INDEX blogs_publish_idx
ON blogs (published_at DESC, score DESC)
WHERE published = TRUE AND status = TRUE;

CREATE INDEX blogs_tags_gin_idx
ON blogs
USING GIN (tags);

CREATE INDEX blogs_search_idx
ON blogs
USING GIN (to_tsvector('english', title));