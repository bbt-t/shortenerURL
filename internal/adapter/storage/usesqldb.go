package storage

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/bbt-t/shortenerURL/internal/entity"
	"github.com/bbt-t/shortenerURL/pkg"

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
	var resultStructs []entity.URLs
	var urlArray []map[string]string

	err := d.db.Select(&resultStructs, "SELECT original_url, short_url FROM items WHERE user_id=$1", userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if len(resultStructs) == 0 {
		return nil, errDBEmpty
	}

	for _, item := range resultStructs {
		temp := make(map[string]string)

		data, _ := json.Marshal(item)
		_ = json.Unmarshal(data, &temp)

		temp["short_url"] = fmt.Sprintf("%s/%s", baseURL, temp["short_url"])
		urlArray = append(urlArray, temp)
	}

	return urlArray, nil
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

func (d *sqlDatabase) DelURLArray(inpJSON []byte, uid string) error {
	inpURLs := pkg.ConvertStrToSlice(string(inpJSON))

	for _, v := range inpURLs {
		_, err := d.db.NamedExec(`UPDATE items SET removed=:removed WHERE id=:id AND user_id=:user_id`,
			map[string]interface{}{
				"removed": true,
				"id":      v,
				"user_id": uid,
			})
		if err != nil {
			return errDBUnknownID
		}
	}
	return nil
}

func (d *sqlDatabase) SaveURLArray(uid uuid.UUID, inpURL []entity.URLBatchInp) error {
	err := saveURLBatch(d.db, uid, inpURL)
	return err
}
