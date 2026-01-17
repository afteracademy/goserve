-- 1. Create Tables
-- ----------------

-- Api Keys Table
CREATE TABLE IF NOT EXISTS api_keys (
    id SERIAL PRIMARY KEY,
    key TEXT NOT NULL,
    permissions TEXT[], -- Using Array type for permissions
    comments TEXT[],    -- Using Array type for comments
    version INTEGER,
    status BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Roles Table
CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    code TEXT NOT NULL UNIQUE,
    status BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Users Table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    status BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Join Table for Users <-> Roles (Many-to-Many relationship)
CREATE TABLE IF NOT EXISTS user_roles (
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    role_id INTEGER REFERENCES roles(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, role_id)
);

-- 2. Insert Data
-- --------------

-- Insert API Key
INSERT INTO api_keys (key, permissions, comments, version, status, created_at, updated_at)
VALUES (
    '1D3F2DD1A5DE725DD4DF1D82BBB37',
    ARRAY['GENERAL'], -- Postgres Array Syntax
    ARRAY['To be used by the xyz vendor'],
    1,
    true,
    NOW(),
    NOW()
);

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
    '$2a$10$psWmSrmtyZYvtIt/FuJL1OLqsK3iR1fZz5.wUYFuSNkkt.EOX9mLa', -- hash of password: changeit
    true, 
    NOW(), 
    NOW()
)
ON CONFLICT (email) DO NOTHING;

-- Map Admin User to ALL Roles
-- This replaces the "db.roles.find({})" logic
INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id
FROM users u
CROSS JOIN roles r
WHERE u.email = 'admin@afteracademy.com'
ON CONFLICT DO NOTHING;