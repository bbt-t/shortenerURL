package sqldb

const _tableItems = `
CREATE TABLE IF NOT EXISTS items (
    id VARCHAR(8), 
    url VARCHAR(512), 
    create_at TIMESTAMP NOT NULL
)`
