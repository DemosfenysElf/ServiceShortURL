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

//go:generate mockgen -source=databaseInterface.go -destination=mocks/mock.go

// StorageInterface
type StorageInterface interface {
	SetURL(ctx context.Context, url string) (short string, err error)
	GetURL(ctx context.Context, short string) (url string, err error)
	Delete(user string, listURL []string)
}

// DatabaseService
type DatabaseService interface {
	Connect(connStr string) error
	CreateTable() error
	Close() error
	Ping(ctx context.Context) error
	GetCount(ctx context.Context) (n int, err error)
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
	connection  *sql.DB
	RandomShort Generator
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
func (db *Database) SetURL(ctx context.Context, url string) (short string, err error) {
	short = db.RandomShort.ShortURL()
	// добавить проверку на оригинальность

	user := GetStructCookies()
	_, err = db.connection.ExecContext(ctx, "insert into ShortenerURL (url,short,nameUser,valueUser,deleted) values ($1,$2,$3,$4,$5)",
		url, short, user.NameUser, user.ValueUser, false)

	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case pgerrcode.UniqueViolation:
			short, _ = db.GetShortURL(ctx, url)
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
func (db *Database) GetURL(ctx context.Context, short string) (url string, err error) {
	deleted := false
	row := db.connection.QueryRowContext(ctx, "select url,deleted from ShortenerURL where short = $1", short)
	err = row.Scan(&url, &deleted)
	fmt.Println(">>>>>URL: ", url, " ; ShortURL: ", short, " ; deleted?: ", deleted)
	if deleted {
		return "", fmt.Errorf("deleted")
	}
	return
}

// GetShortURL передаём оригинальный URL
// получаем сохраненный в БД короткий URL
func (db *Database) GetShortURL(ctx context.Context, url string) (short string, err error) {
	row := db.connection.QueryRowContext(ctx, "select short from ShortenerURL where url = $1", url)
	err = row.Scan(&short)
	return
}

// Delete передаём данные о пользователе и список URL
// удаляем все URL из списка, принадлежащие этому пользователю
func (db *Database) Delete(user string, listURL []string) {
	fmt.Println(">>>BD_Delete_list<<<  ", listURL, "User: ", user)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*11)
	defer cancel()
	for _, u := range listURL {
		_, err := db.connection.ExecContext(ctx, "UPDATE ShortenerURL SET deleted = true WHERE short=$1 AND valueUser=$2", u, user)
		if err != nil {
			fmt.Println(">>>>>>>>>>>>>>>>>>>>>", err)
		}
	}
}

// SetConnection для тестирования с помощью mock
func (db *Database) SetConnection(conn *sql.DB) {
	db.connection = conn
}

// GetCount возвращаем количество строк в таблице
func (db *Database) GetCount(ctx context.Context) (n int, err error) {
	row := db.connection.QueryRowContext(ctx, "select count(url) from ShortenerURL")
	err = row.Scan(&n)
	return
}
