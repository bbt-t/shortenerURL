package storage

const _tableItems = `
CREATE TABLE IF NOT EXISTS users (
	user_id UUID NOT NULL PRIMARY KEY,
    create_at TIMESTAMP NOT NULL
);
CREATE TABLE IF NOT EXISTS items (
	user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
	original_url VARCHAR(512),
	short_url VARCHAR(32),
    create_at TIMESTAMP NOT NULL
);
`
