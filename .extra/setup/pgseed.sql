-- Enable UUID
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- 1. Create Tables
-- ----------------

-- Api Keys Table
CREATE TABLE IF NOT EXISTS api_keys (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    key TEXT NOT NULL UNIQUE,
    permissions TEXT[],
    comments TEXT[],
    version INTEGER,
    status BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Api Keys Indexes
CREATE INDEX IF NOT EXISTS api_keys_key_status_idx
ON api_keys (key, status);

-- Roles Table
CREATE TABLE IF NOT EXISTS roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code TEXT NOT NULL UNIQUE,
    status BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Users Table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
		profile_pic_url TEXT,
		verified BOOLEAN DEFAULT FALSE,
    status BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Join Table for Users <-> Roles
CREATE TABLE IF NOT EXISTS user_roles (
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    role_id UUID REFERENCES roles(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, role_id)
);

-- Keystore Table
CREATE TABLE IF NOT EXISTS keystore (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	p_key TEXT NOT NULL,
	s_key TEXT NOT NULL,
	status BOOLEAN DEFAULT TRUE,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Keystore Table Indexes
CREATE INDEX IF NOT EXISTS keystore_user_status_idx
ON keystore (user_id, status);

CREATE INDEX IF NOT EXISTS keystore_user_pkey_status_idx
ON keystore (user_id, p_key, status);

CREATE INDEX IF NOT EXISTS keystore_user_pkey_skey_status_idx
ON keystore (user_id, p_key, s_key, status);

-- Messages Table
CREATE TABLE messages (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	type TEXT NOT NULL,
	msg TEXT NOT NULL,
	status BOOLEAN DEFAULT TRUE,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Blogs Table
CREATE TABLE IF NOT EXISTS blogs (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	title TEXT NOT NULL,
	description TEXT NOT NULL,
	text TEXT,
	draft_text TEXT NOT NULL,
	tags TEXT[],
	author_id UUID NOT NULL REFERENCES users(id),
	img_url TEXT,
	slug TEXT NOT NULL UNIQUE,
	score DOUBLE PRECISION DEFAULT 0.01,
	submitted BOOLEAN DEFAULT FALSE,
	drafted BOOLEAN DEFAULT TRUE,
	published BOOLEAN DEFAULT FALSE,
	status BOOLEAN DEFAULT TRUE,
	published_at TIMESTAMP,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Blogs Table Indexes
CREATE INDEX IF NOT EXISTS blogs_publish_idx
ON blogs (published_at DESC, score DESC)
WHERE published = TRUE AND status = TRUE;

CREATE INDEX IF NOT EXISTS blogs_tags_gin_idx
ON blogs
USING GIN (tags);

CREATE INDEX IF NOT EXISTS blogs_search_idx
ON blogs
USING GIN (to_tsvector('english', title));

-- 2. Insert Data
-- --------------

-- Insert API Key
INSERT INTO api_keys (key, permissions, comments, version, status, created_at, updated_at)
VALUES (
    '1D3F2DD1A5DE725DD4DF1D82BBB37',
    ARRAY['GENERAL'],
    ARRAY['To be used by the xyz vendor'],
    1,
    true,
    NOW(),
    NOW()
)
ON CONFLICT (key) DO NOTHING;

-- Insert Roles
INSERT INTO roles (code, status, created_at, updated_at)
VALUES 
    ('LEARNER', true, NOW(), NOW()),
    ('AUTHOR', true, NOW(), NOW()),
    ('EDITOR', true, NOW(), NOW()),
    ('ADMIN', true, NOW(), NOW())
ON CONFLICT (code) DO NOTHING;

-- Insert Admin User
INSERT INTO users (name, email, password, status, created_at, updated_at)
VALUES (
    'Admin', 
    'admin@afteracademy.com', 
    '$2a$10$psWmSrmtyZYvtIt/FuJL1OLqsK3iR1fZz5.wUYFuSNkkt.EOX9mLa',
    true, 
    NOW(), 
    NOW()
)
ON CONFLICT (email) DO NOTHING;

-- Map Admin User to ALL Roles
INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id
FROM users u
CROSS JOIN roles r
WHERE u.email = 'admin@afteracademy.com'
ON CONFLICT DO NOTHING;