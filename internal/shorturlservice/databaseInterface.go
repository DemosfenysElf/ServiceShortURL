package shorturlservice

import (
	"context"
	"database/sql"
	"time"
)

type DatabaseInterface interface {
	Connect(connStr string) error
	Close() error
	Ping() error
}

var stringShortenerURL string = `CREATE TABLE ShortenerURL(
url            varchar(64),
short          varchar(32),
nameAut        varchar(32),
valueAut       varchar(32)
)`

type Database struct {
	connection *sql.DB
	ctx        context.Context
}

func InitDB() (*Database, error) {
	return &Database{ctx: context.Background()}, nil
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
	err = db.Ping()
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) Close() error {
	return db.connection.Close()
}

func (db *Database) Ping() error {
	ctx, cancel := context.WithTimeout(db.ctx, 1*time.Second)
	defer cancel()
	if err := db.connection.PingContext(ctx); err != nil {
		return err
	}
	return nil
}

func (db *Database) SetURL(url string) (short string) {
	short = shortURL()
	// добавить проверку на оригинальность
	//ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()

	user := GetStructCookies()
	_, err := db.connection.Exec("insert into ShortenerURL(url,short,nameAut,valueAut) values ($1,$2,$3,$4)", url, short, user.NameUser, user.ValueUser)
	if err != nil {
		return ""
	}
	return short
}

func (db *Database) GetURL(short string) (url string, err error) {
	row := db.connection.QueryRow("select url from ShortenerURL where short = $1", short)
	err = row.Scan(&url)
	return
}

func (db *Database) CreateTable() error {
	_, err := db.connection.Exec(stringShortenerURL)
	return err
}