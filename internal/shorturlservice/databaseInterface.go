package shorturlservice

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgerrcode"
	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// DatabaseService
type DatabaseService interface {
	Connect(connStr string) error
	Close() error
	Ping(ctx context.Context) error
}

var stringShortenerURL = `CREATE TABLE ShortenerURL(
url            varchar(64),
short          varchar(32),
nameUser        varchar(32),
valueUser       varchar(32),
deleted			bool
)`

// Database connection *sql.DB
type Database struct {
	connection *sql.DB
}

// Подключние к БД по пути
func (db *Database) Connect(connStr string) (err error) {
	db.connection, err = sql.Open("pgx", connStr)
	if err != nil {
		return err
	}

	db.CreateTable()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err = db.Ping(ctx)
	if err != nil {
		return err
	}
	return nil
}

// CreateTable() создание таблицы
func (db *Database) CreateTable() error {
	_, err := db.connection.Exec(stringShortenerURL)
	if err != nil {
		return err
	}
	_, err = db.connection.Exec("CREATE UNIQUE INDEX URL_index ON ShortenerURL (url)")
	return err
}

// Close закрывашка
func (db *Database) Close() error {
	return db.connection.Close()
}

// Ping проверяет, что соединение с базой данных все еще работает
func (db *Database) Ping(ctx context.Context) error {
	if err := db.connection.PingContext(ctx); err != nil {
		return err
	}
	return nil
}

// SetURL передаём оригинальный URL
// получаем сгенерированный короткий URL
// вместе с данными о пользователе сохраняем в БД
func (db *Database) SetURL(url string) (short string, err error) {
	short = shortURL()
	// добавить проверку на оригинальность

	user := GetStructCookies()
	_, err = db.connection.Exec("insert into ShortenerURL (url,short,nameUser,valueUser,deleted) values ($1,$2,$3,$4,$5)",
		url, short, user.NameUser, user.ValueUser, false)

	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case pgerrcode.UniqueViolation:
			short, _ = db.GetShortURL(url)
			return short, err
		default:
			return "", err
		}

	}
	return short, err
}

// GetURL передаём короткий URL
// вместе с данными о пользователе сохраняем в БД
// получаем сгенерированный оригинальный URL
func (db *Database) GetURL(short string) (url string, err error) {
	deleted := false
	row := db.connection.QueryRow("select url,deleted from ShortenerURL where short = $1", short)
	err = row.Scan(&url, &deleted)
	fmt.Println(">>>>>URL: ", url, " ; shortURL: ", short, " ; deleted?: ", deleted)
	if deleted {
		return "", fmt.Errorf("deleted")
	}
	return
}

// GetShortURL передаём оригинальный URL
// получаем сохраненный в БД короткий URL
func (db *Database) GetShortURL(url string) (short string, err error) {
	row := db.connection.QueryRow("select short from ShortenerURL where url = $1", url)
	err = row.Scan(&short)
	return
}

// Delete передаём данные о пользователе и список URL
// удаляем все URL из списка, принадлежащие этому пользователю
func (db *Database) Delete(user string, listURL []string) {
	fmt.Println(">>>BD_Delete_list<<<  ", listURL, "User: ", user)
	//defer wg.Done()
	for _, u := range listURL {
		_, err := db.connection.Exec("UPDATE ShortenerURL SET deleted = true WHERE short=$1 AND valueUser=$2", u, user)
		if err != nil {
			fmt.Println(">>>>>>>>>>>>>>>>>>>>>", err)
		}
	}

}
