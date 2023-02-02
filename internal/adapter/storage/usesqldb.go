package storage

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"

	"github.com/bbt-t/shortenerURL/internal/entity"

	"github.com/gofrs/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type sqlDatabase struct {
	db *pgxpool.Pool
}

func NewSQLDatabase(dsn string) DatabaseRepository {
	/*
		Selects DB and initializing. Create tables.
		param nameDB: received parameter (flag) to select db
		return: db-object or nil
	*/
	ctx := context.Background()
	db, err := pgxpool.New(ctx, dsn)

	if err != nil {
		log.Println(err)
		return nil
	}

	createTable(ctx, db, _tableItems /* SQL command */)
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
	ctx := context.Background()
	err := d.db.Ping(ctx)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Postgres is READY!")
	}
	return err
}

func (d *sqlDatabase) DelURLArray(ctx context.Context, userID uuid.UUID, inpURLs []string) error {
	err := deleteURLArray(ctx, d.db, userID, inpURLs)
	return err
}

func (d *sqlDatabase) SaveURLArray(ctx context.Context, uid uuid.UUID, inpURL []entity.URLBatchInp) error {
	err := saveURLBatch(ctx, d.db, uid, inpURL)
	return err
}
