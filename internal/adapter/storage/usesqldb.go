package storage

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type sqlDatabase struct {
	db *sqlx.DB
}

type URLs struct {
	OriginalURL string `db:"original_url" json:"original_url"`
	ShortURL    string `db:"short_url" json:"short_url"`
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
	var resultStructs []URLs
	var urlArray []map[string]string

	err := d.db.Select(&resultStructs, "SELECT original_url, short_url FROM items WHERE user_id=$1", userID)
	if err != nil {
		fmt.Println(err)
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

//func (d *sqlDatabase) DelURLArray(inputURLJSON []byte, UID string) error {
//	var idu int64
//
//	ctx := context.Background()
//	vURL := strings.ReplaceAll(string(inputURLJSON), " ", "")
//	vURL = strings.ReplaceAll(strings.ReplaceAll(vURL, "[", ""), "]", "")
//
//	valURL := strings.Split(strings.ReplaceAll(vURL, "\"", ""), ",")
//	fmt.Println("Split url short ", valURL)
//	if len(valURL) > 20 {
//		// batch
//		batch := &pgx.Batch{}
//		for _, v := range valURL {
//			batch.Queue("UPDATE items SET deleted=true "+
//				"where url_id=(select id from url where url_short=$1) and user_id=$2", v, idu)
//		}
//		br := m.db.SendBatch(context.Background(), batch)
//
//		ct, err := br.Exec()
//		if err != nil {
//			fmt.Println("Not Updated users_url(user_id,url_id) ", err)
//		}
//		if ct.RowsAffected() != 1 {
//			fmt.Println("ct.RowsAffected()", ct.RowsAffected())
//		}
//		br.Close()
//
//	} else {
//		// обновляем одним запросом, списком
//		vURL := strings.ReplaceAll(vURL, "\"", "'")
//		query := "UPDATE users_url set deleted=true where " +
//			"url_id in (select id from url where url_short in (" + vURL +
//			") ) and user_id=$1"
//		fmt.Println("query =", query)
//		if _, err := m.db.Exec(ctx, query, idu); err != nil {
//			fmt.Println("Not Updated users_url(user_id,url_id) ", err)
//			return err
//		}
//	}
//	return nil
//}
