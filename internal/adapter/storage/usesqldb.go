package storage

import (
	"encoding/json"
	"log"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type sqlDatabase struct {
	db *sqlx.DB
}

type User struct {
	HashURL     string `db:"short_url"`
	OriginalURL string `db:"original_url"`
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
		//log.Println(errDBNotSelected)
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
	result, err := getHashURL(d.db, k /* id (hashed url) */)
	if err != nil {
		return "", err
	}
	return result, nil
}

func (d *sqlDatabase) GetURLArrayByUser(userID uuid.UUID) ([]map[string]string, error) {
	userMap := User{}
	allURL := map[string]string{}

	if err := d.db.Get(
		&userMap,
		"SELECT short_url, original_url FROM items WHERE user_id=$1",
		userID,
	); err != nil {
		return nil, err
	}

	//inMap := structs.Map(result)

	data, errMJson := json.Marshal(userMap)
	if errMJson != nil {
		log.Println(errMJson)
		return nil, errMJson
	}
	if errUJson := json.Unmarshal(data, &allURL); errUJson != nil {
		log.Println(errUJson)
		return nil, errUJson
	}

	result := convertToArrayMap(allURL)

	return result, nil
}

func (d *sqlDatabase) SaveShortURL(userID uuid.UUID, k, v string) error {
	/*
		Adding info to the DB.
	*/
	err := saveURL(d.db, userID, k /* hashed url */, v /* original url */)
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
