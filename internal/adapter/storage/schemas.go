package storage

const _tableItems = `
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS users (
	user_id UUID PRIMARY KEY,
	create_at TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE TABLE IF NOT EXISTS items (
	id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
	user_id UUID NOT NULL REFERENCES users(user_id),
	original_url VARCHAR(512),
	short_url VARCHAR(32),
	create_at TIMESTAMP NOT NULL DEFAULT NOW(),
	deleted bool not null DEFAULT false
	
);
CREATE INDEX IF NOT EXISTS idx_user_id ON users (user_id);
CREATE INDEX IF NOT EXISTS idx_url_id ON items (id);
`
