package storage

import (
	"context"
	"log"

	"github.com/bbt-t/shortenerURL/internal/entity"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type sqlDatabase struct {
	db *sqlx.DB
}

func NewSQLDatabase(dsn string) DatabaseRepository {
	/*
		Selects DB and initializing. Create tables.
		param nameDB: received parameter (flag) to select db
		return: db-object or nil
	*/
	db, err := sqlx.Connect("postgres", dsn)

	if err != nil {
		log.Println(err)
		return nil
	}

	createTable(db, _tableItems /* SQL command */)
	return &sqlDatabase{
		db: db,
	}
}

func (d *sqlDatabase) NewUser(userID uuid.UUID) {
	/*
		Adds new user.
	*/
	if checkUser(d.db, userID) {
		return
	}
	addNewUser(d.db, userID)
}

func (d *sqlDatabase) GetOriginalURL(k string) (string, error) {
	/*
		Search for info by ID.
		param k: id by which we search in the DB
	*/
	result, err := getOriginalURL(d.db, k /* id (hashed url) */)
	if err != nil {
		return "", err
	}
	return result, nil
}

func (d *sqlDatabase) GetURLArrayByUser(userID uuid.UUID, baseURL string) ([]map[string]string, error) {
	/*
		Gets all pairs "original" - "short" urls previously saved by the user.
	*/
	result, err := getOriginalURLArray(d.db, userID, baseURL)
	return result, err
}

func (d *sqlDatabase) SaveShortURL(userID uuid.UUID, shortURL, originalURL string) error {
	/*
		Adding info to the DB.
	*/
	err := saveURL(d.db, userID, shortURL, originalURL)
	return err
}

func (d *sqlDatabase) PingDB() error {
	/*
		Checking connection with ctx.Background.
	*/
	err := d.db.Ping()
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Postgres is READY!")
	}
	return err
}

func (d *sqlDatabase) DelURLArray(ctx context.Context, uid uuid.UUID, inpJSON []byte) error {
	//err := deleteURLArray(ctx, d.db, uid, inpJSON)
	err := deleteURLArrayQueue(ctx, d.db, uid, inpJSON)
	return err
}

func (d *sqlDatabase) SaveURLArray(ctx context.Context, uid uuid.UUID, inpURL []entity.URLBatchInp) error {
	err := saveURLBatch(ctx, d.db, uid, inpURL)
	return err
}
