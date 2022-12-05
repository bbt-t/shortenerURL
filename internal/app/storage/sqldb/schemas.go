package sqldb

import "time"

const tableItems = `
CREATE TABLE IF NOT EXISTS items (
    id VARCHAR(8), 
    url VARCHAR(512), 
    create_at TIMESTAMP NOT NULL
)`

type items struct {
	Id       string    `db:"id"`
	Url      string    `db:"url"`
	CreateAt time.Time `db:"create_at"`
}
