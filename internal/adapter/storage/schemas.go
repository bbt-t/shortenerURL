package storage

const _tableItems = `
CREATE TABLE IF NOT EXISTS items (
    user_id UUID,
    short_url VARCHAR(32),
    original_url VARCHAR(512), 
    create_at TIMESTAMP NOT NULL
)`

const _index = `
CREATE UNIQUE INDEX IF NOT EXISTS short_url ON items(short_url)
`
