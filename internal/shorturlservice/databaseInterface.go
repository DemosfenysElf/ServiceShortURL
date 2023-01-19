package shorturlservice

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgerrcode"
)

type DatabaseService interface {
	Connect(connStr string) error
	Close() error
	Ping(ctx context.Context) error
	DeleteURL(user string, listURL []string) error
}

var stringShortenerURL = `CREATE TABLE ShortenerURL(
url            varchar(64),
short          varchar(32),
nameUser        varchar(32),
valueUser       varchar(32),
deleted			bool
)`

type Database struct {
	connection *sql.DB
}

func InitDB() (*Database, error) {
	return &Database{}, nil
}

func (db *Database) Connect(connStr string) (err error) {
	db.connection, err = sql.Open("pgx", connStr)
	if err != nil {
		return err
	}
	db.connection.Exec("Drop TABLE ShortenerURL")
	err = db.CreateTable()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err = db.Ping(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) Close() error {
	return db.connection.Close()
}

func (db *Database) Ping(ctx context.Context) error {
	if err := db.connection.PingContext(ctx); err != nil {
		return err
	}
	return nil
}

func (db *Database) SetURL(url string) (short string, err error) {
	fmt.Println(">>>>>>>>>SetURL, DB")
	short = shortURL()
	// добавить проверку на оригинальность

	user := GetStructCookies()
	_, err = db.connection.Exec("insert into ShortenerURL (url,short,nameUser,valueUser,deleted) values ($1,$2,$3,$4,$5)",
		url, short, user.NameUser, user.ValueUser, false)

	var sErr string
	if err != nil {
		sErr = err.Error()
		if strings.Contains(sErr, pgerrcode.UniqueViolation) {
			short, _ = db.GetShortURL(url)
			return short, err
		}
		return "", err
	}

	return short, err
}

func (db *Database) GetURL(short string) (url string, err error) {
	deleted := false
	row := db.connection.QueryRow("select url,deleted from ShortenerURL where short = $1", short)
	err = row.Scan(&url, &deleted)
	fmt.Println(">>>>>URL: ", url)
	fmt.Println(">>>>>deleted: ", deleted)
	if deleted {
		return "", fmt.Errorf("deleted")
	}
	return
}

func (db *Database) GetShortURL(url string) (short string, err error) {
	row := db.connection.QueryRow("select short from ShortenerURL where url = $1", url)
	err = row.Scan(&short)
	return
}

func (db *Database) CreateTable() error {
	_, err := db.connection.Exec(stringShortenerURL)
	if err != nil {
		return err
	}
	_, err = db.connection.Exec("CREATE UNIQUE INDEX URL_index ON ShortenerURL (url)")
	return err
}

func (db *Database) DeleteURL(user string, listURL []string) error {
	fmt.Println(">>>Delete_list<<<  ", listURL, "User: ", user)

	for _, u := range listURL {
		_, err := db.connection.Exec("UPDATE ShortenerURL SET deleted = true WHERE (url=$1)&&(valueUser=$2)", u, user)
		if err != nil {
			return err
		}
	}
	return nil
}
