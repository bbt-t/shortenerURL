package storage

import (
	"encoding/json"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
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

func (d *sqlDatabase) SaveShortURL(userID uuid.UUID, k, v string) error {
	/*
		Adding info to the DB.
	*/
	err := saveURl(d.db, userID, k /* hashed url */, v /* original url */)
	return err
}

func (d *sqlDatabase) GetOriginalURL(k string) (string, error) {
	/*
		Search for info by ID.
		param k: id by which we search in the DB
	*/
	result, err := getHashUrl(d.db, k /* id (hashed url) */)
	if err != nil {
		return "", err
	}
	return result, nil
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

func (d *sqlDatabase) NewUser(userID uuid.UUID) {
	if checkUser(d.db, userID) {
		return
	}
	addNewUser(d.db, userID)
}

func (d *sqlDatabase) GetURLArrayByUser(userID uuid.UUID) (map[string]string, error) {
	result := User{}
	inMap := map[string]string{}

	err := d.db.Get(&result, "SELECT short_url, original_url FROM items WHERE user_id=$1", userID)

	//inMap := structs.Map(result)

	data, errMJson := json.Marshal(result)
	if err != nil {
		log.Println(errMJson)
		return nil, errMJson
	}
	if errUJson := json.Unmarshal(data, &inMap); err != nil {
		log.Println(errUJson)
		return nil, errUJson
	}

	return inMap, err
}
