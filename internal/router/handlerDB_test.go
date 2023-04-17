package router

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/labstack/echo"

	"ServiceShortURL/internal/shorturlservice"
)

func TestPostDB(t *testing.T) {
	type want struct {
		codePost int
		response string
	}
	tests := []struct {
		name           string
		want           want
		wantErr        bool
		generatorShort *shorturlservice.TestGenerator
		generatorUser  *shorturlservice.TestGeneratorUser
		urlResult      []string
		urllist        []string
		bodyURLJSON    string
		result         driver.Result
	}{
		{
			name: "TestBDApiShortenBatch1",
			want: want{
				codePost: 201,
				response: `{"status":"ok"}`,
			},
			generatorShort: &shorturlservice.TestGenerator{Result: []string{"Short1", "Short2", "Short3", "Short4"}, Index: 0},
			generatorUser:  &shorturlservice.TestGeneratorUser{Result: []byte("UserUserUserUser")},
			urlResult: []string{
				"/Short1",
				"/Short2",
				"/Short3",
				"/Short4",
			},
			urllist: []string{
				"https://www.youtube.com/watch?v=UK7yzgVpnDA",
				"https://metanit.com/",
				"https://otion.site/context-19ec3dbb0cd4ec",
				"https://stepik.org/lesson/359395",
			},

			bodyURLJSON: `[{"correlation_id":"baseurl1","original_url":"https://www.youtube.com/watch?v=UK7yzgVpnDA"},
							{"correlation_id":"baseurl2","original_url":"https://metanit.com/"},
							{"correlation_id":"baseurl3","original_url":"https://otion.site/context-19ec3dbb0cd4ec"},
							{"correlation_id":"baseurl4","original_url":"https://stepik.org/lesson/359395"}]`,
			result: sqlmock.NewResult(1, 1),
		},
		{
			name: "TestBDApiShortenBatch2",
			want: want{
				codePost: 500,
				response: `{"status":"ok"}`,
			},
			generatorShort: &shorturlservice.TestGenerator{Result: []string{"Short1", "Short2", "Short3", "Short4"}, Index: 0},
			generatorUser:  &shorturlservice.TestGeneratorUser{Result: []byte("UserUserUserUser")},
			urlResult: []string{
				"/Short1",
				"/Short2",
				"/Short3",
				"/Short4",
			},
			urllist: []string{
				"https://www.youtube.com/watch?v=UK7yzgVpnDA",
				"https://metanit.com/",
				"https://otion.site/context-19ec3dbb0cd4ec",
				"https://stepik.org/lesson/359395",
			},

			bodyURLJSON: `[{"correlation_id":"baseurl1","original_url":"https://www.youtube.com/watch?v=UK7yzgVpnDA"},
							{"correlation_id":"baseurl2","original_url":"https://metanit.com/"},
							{"correlation_id":"baseurl3","original_url":"https://otion.site/context-19ec3dbb0cd4ec"},
							{"correlation_id":"baseurl4","original_url":"https://stepik.org/lesson/359395"}`,
			result: sqlmock.NewResult(1, 1),
		},
		{
			name: "TestBDApiShortenBatch3",
			want: want{
				codePost: 201,
				response: `{"status":"ok"}`,
			},
			generatorShort: &shorturlservice.TestGenerator{Result: []string{"Short1", "Short2", "Short3", "Short4"}, Index: 0},
			generatorUser:  &shorturlservice.TestGeneratorUser{Result: []byte("UserUserUserUser")},
			urlResult: []string{
				"/Short1",
				"/Short2",
				"/Short3",
				"/Short4",
			},
			urllist: []string{
				"https://www.youtube.com/watch?v=UK7yzgVpnDA",
				"https://metanit.com/",
				"https://otion.site/context-19ec3dbb0cd4ec",
				"https://stepik.org/lesson/359395",
			},

			bodyURLJSON: `[{"correlation_id":"baseurl1","original_url":"https://www.youtube.com/watch?v=UK7yzgVpnDA"},
							{"correlation_id":"baseurl2","original_url":"https://metanit.com/"},
							{"correlation_id":"baseurl3","original_url":"https://metanit.com/"},
							{"correlation_id":"baseurl4","original_url":"https://stepik.org/lesson/359395"}]`,
			result: sqlmock.NewResult(1, 1),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()
			rout := InitTestServer()
			rout.Cfg.Storage = testStorageURL
			mockDB := &shorturlservice.Database{RandomShort: tt.generatorShort}

			mockDB.SetConnection(db)
			rout.DB = mockDB
			rout.StorageInterface = mockDB
			rout.GeneratorUsers = tt.generatorUser
			user := rout.kykiDB()
			//
			for i := range tt.generatorShort.Result {
				mock.ExpectExec("insert into ShortenerURL").
					WithArgs(tt.urllist[i], tt.generatorShort.ShortURL(), user.Name, user.Value, false).
					WillReturnResult(tt.result)
			}
			tt.generatorShort.Index = 0
			//

			e := echo.New()
			request := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", strings.NewReader(tt.bodyURLJSON))

			request.AddCookie(user)
			responseRecorder := httptest.NewRecorder()
			c := e.NewContext(request, responseRecorder)

			rout.Serv = e
			rout.PostAPIShortenBatch(c)
			response := responseRecorder.Result()

			defer response.Body.Close()
			resBody, err := io.ReadAll(response.Body)
			if err != nil {
				t.Fatal(err)
			}

			if response.StatusCode != tt.want.codePost {
				t.Errorf("Expected status code %d, got %d", tt.want.codePost, responseRecorder.Code)
			}

			if response.StatusCode == http.StatusCreated {

				urlBatch := []shortURLApiShortenBatch{}
				err = json.Unmarshal(resBody, &urlBatch)
				if err != nil {
					t.Fatal(err)
				}
				for i, batch := range urlBatch {
					if batch.ShortURL != tt.urlResult[i] {
						t.Errorf("Короткая ссылка не соответствует ожидаемой")
					}
					fmt.Println(batch)
				}

			}

		})
	}
}

func TestDBGet(t *testing.T) {
	type want struct {
		codeGet  int
		url      string
		response string
	}
	tests := []struct {
		name           string
		want           want
		wantErr        bool
		generatorShort *shorturlservice.TestGenerator
		generatorUser  *shorturlservice.TestGeneratorUser
		urlShort       string
		bodyURL        string
		del            bool
	}{
		{
			name: "TestDBGet1",
			want: want{
				codeGet:  307,
				url:      "https://www.youtube.com/watch?v=UK7yzgVpnDA",
				response: `{"status":"ok"}`,
			},
			generatorShort: &shorturlservice.TestGenerator{Result: []string{"Short1"}, Index: 0},
			generatorUser:  &shorturlservice.TestGeneratorUser{Result: []byte("UserUserUserUser")},
			urlShort:       "Short1",
			bodyURL:        "https://www.youtube.com/watch?v=UK7yzgVpnDA",
			del:            false,
		},
		{
			name: "TestDBGet2",
			want: want{
				codeGet:  410,
				url:      "https://www.youtube.com/watch?v=UK7yzgVpnDA",
				response: `{"status":"ok"}`,
			},
			generatorShort: &shorturlservice.TestGenerator{Result: []string{"Short2"}, Index: 0},
			generatorUser:  &shorturlservice.TestGeneratorUser{Result: []byte("UserUserUserUser")},
			urlShort:       "Short2",
			bodyURL:        "https://www.youtube.com/watch?v=UK7yzgVpnDA",
			del:            true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()
			rout := InitTestServer()
			rout.Cfg.Storage = testStorageURL
			mockDB := &shorturlservice.Database{RandomShort: tt.generatorShort}

			mockDB.SetConnection(db)
			rout.DB = mockDB
			rout.StorageInterface = mockDB
			rout.GeneratorUsers = tt.generatorUser
			user := rout.kykiDB()
			//
			row := sqlmock.NewRows([]string{"url", "deleted"}).
				AddRow(tt.bodyURL, tt.del)

			mock.ExpectQuery("select url,deleted from ShortenerURL where ?").
				WithArgs(tt.urlShort).WillReturnRows(row)

			tt.generatorShort.Index = 0
			//

			e := echo.New()
			request := httptest.NewRequest(http.MethodPost, "/"+tt.urlShort, nil)

			request.AddCookie(user)
			responseRecorder := httptest.NewRecorder()
			c := e.NewContext(request, responseRecorder)

			rout.Serv = e
			rout.GetShortToURL(c)
			response := responseRecorder.Result()

			resBody := response.Header.Get("Location")
			defer response.Body.Close()

			if response.StatusCode != tt.want.codeGet {
				t.Errorf("Expected status code %d, got %d", tt.want.codeGet, responseRecorder.Code)
			}

			if response.StatusCode == http.StatusTemporaryRedirect {

				if resBody != tt.want.url {
					t.Errorf("Ожидаемое значение %s, полученное %s", tt.want.url, resBody)
				}

			}

		})
	}
}

func TestDBDelete(t *testing.T) {
	type want struct {
		codeDel  int
		response string
	}
	tests := []struct {
		name           string
		want           want
		wantErr        bool
		generatorShort *shorturlservice.TestGenerator
		generatorUser  *shorturlservice.TestGeneratorUser
		bodyURLJSON    string
	}{
		{
			name: "TestBDDel",
			want: want{
				codeDel:  202,
				response: `{"status":"ok"}`,
			},
			generatorShort: &shorturlservice.TestGenerator{Result: []string{"Short1", "Short2", "Short3", "Short4"}, Index: 0},
			generatorUser:  &shorturlservice.TestGeneratorUser{Result: []byte("UserUserUserUser")},

			bodyURLJSON: `["/Short1","/Short2","/Short3","/Short4"]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}

			rout := InitTestServer()
			rout.Cfg.Storage = testStorageURL
			mockDB := &shorturlservice.Database{RandomShort: tt.generatorShort}

			mockDB.SetConnection(db)
			rout.DB = mockDB
			rout.StorageInterface = mockDB
			rout.GeneratorUsers = tt.generatorUser
			user := rout.kykiDB()
			//
			for range tt.generatorShort.Result {
				mock.ExpectExec("UPDATE ShortenerURL SET deleted = true WHERE short=? AND valueUser=?").
					WithArgs(tt.generatorShort.ShortURL(), user.Value).
					WillReturnResult(sqlmock.NewResult(1, 1))
			}
			tt.generatorShort.Index = 0
			//

			e := echo.New()
			request := httptest.NewRequest(http.MethodDelete, "/api/user/urls", strings.NewReader(tt.bodyURLJSON))

			request.AddCookie(user)
			responseRecorder := httptest.NewRecorder()
			c := e.NewContext(request, responseRecorder)

			rout.Serv = e
			rout.DeleteAPIUserURLs(c)
			response := responseRecorder.Result()
			defer response.Body.Close()

			if response.StatusCode != tt.want.codeDel {
				t.Errorf("Expected status code %d, got %d", tt.want.codeDel, responseRecorder.Code)
			}

		})
	}
}

func TestDBPostUniqueViolation(t *testing.T) {
	type want struct {
		codePost int
		response string
	}
	tests := []struct {
		name           string
		want           want
		wantErr        bool
		generatorShort *shorturlservice.TestGenerator
		generatorUser  *shorturlservice.TestGeneratorUser
		urlResult      string
		urllist        string
		bodyURL        string
		errUnique      *pgconn.PgError
	}{
		{
			name: "TestBDApiShortenBatch1",
			want: want{
				codePost: http.StatusConflict,
				response: `{"status":"ok"}`,
			},
			generatorShort: &shorturlservice.TestGenerator{Result: []string{"Short1"}, Index: 0},
			generatorUser:  &shorturlservice.TestGeneratorUser{Result: []byte("UserUserUserUser")},
			urlResult:      "Short1",
			urllist:        "https://www.youtube.com/watch?v=UK7yzgVpnDA",
			bodyURL:        "https://www.youtube.com/watch?v=UK7yzgVpnDA",
			errUnique:      &pgconn.PgError{Code: pgerrcode.UniqueViolation},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()
			rout := InitTestServer()
			rout.Cfg.Storage = testStorageURL
			mockDB := &shorturlservice.Database{RandomShort: tt.generatorShort}

			mockDB.SetConnection(db)
			rout.DB = mockDB
			rout.StorageInterface = mockDB
			rout.GeneratorUsers = tt.generatorUser
			user := rout.kykiDB()
			//
			mock.ExpectExec("insert into ShortenerURL").
				WithArgs(tt.urllist, tt.generatorShort.ShortURL(), user.Name, user.Value, false).
				WillReturnError(tt.errUnique)

			//db.connection.QueryRowContext(ctx, "select short from ShortenerURL where url = $1", url)
			row := sqlmock.NewRows([]string{"short"}).
				AddRow(tt.urlResult)

			mock.ExpectQuery("select short from ShortenerURL where ?").
				WithArgs(tt.urllist).
				WillReturnRows(row)

			tt.generatorShort.Index = 0
			//

			e := echo.New()
			request := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", strings.NewReader(tt.bodyURL))

			request.AddCookie(user)
			responseRecorder := httptest.NewRecorder()
			c := e.NewContext(request, responseRecorder)

			rout.Serv = e
			rout.PostURLToShort(c)
			response := responseRecorder.Result()

			defer response.Body.Close()
			resBody, err := io.ReadAll(response.Body)
			if err != nil {
				t.Fatal(err)
			}

			if response.StatusCode != tt.want.codePost {
				t.Errorf("Expected status code %d, got %d", tt.want.codePost, responseRecorder.Code)
			}

			if string(resBody[1:]) != tt.urlResult {
				t.Errorf("Короткая ссылка не соответствует ожидаемой")
			}
			fmt.Println(string(resBody))

		})
	}
}

func TestDBPing(t *testing.T) {
	type want struct {
		codeGet  int
		response string
	}
	tests := []struct {
		name           string
		want           want
		generatorShort *shorturlservice.TestGenerator
		generatorUser  *shorturlservice.TestGeneratorUser
		errPing        error
	}{
		{
			name: "TestDBGetPing1",
			want: want{
				codeGet:  http.StatusOK,
				response: `{"status":"ok"}`,
			},
			generatorShort: &shorturlservice.TestGenerator{Result: []string{"Short1"}, Index: 0},
			generatorUser:  &shorturlservice.TestGeneratorUser{Result: []byte("UserUserUserUser")},
			errPing:        nil,
		},
		{
			name: "TestDBGetPing2",
			want: want{
				codeGet:  http.StatusInternalServerError,
				response: `{"status":"ok"}`,
			},
			generatorShort: &shorturlservice.TestGenerator{Result: []string{"Short1"}, Index: 0},
			generatorUser:  &shorturlservice.TestGeneratorUser{Result: []byte("UserUserUserUser")},
			errPing:        errors.New("ping fail"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()
			rout := InitTestServer()
			rout.Cfg.Storage = testStorageURL
			mockDB := &shorturlservice.Database{RandomShort: tt.generatorShort}

			mockDB.SetConnection(db)
			rout.DB = mockDB
			rout.StorageInterface = mockDB
			rout.GeneratorUsers = tt.generatorUser
			user := rout.kykiDB()
			//
			mock.ExpectPing().WillReturnError(tt.errPing)
			//
			e := echo.New()
			request := httptest.NewRequest(http.MethodGet, "/ping", nil)

			request.AddCookie(user)
			responseRecorder := httptest.NewRecorder()
			c := e.NewContext(request, responseRecorder)

			rout.Serv = e
			rout.GetPingDB(c)
			response := responseRecorder.Result()
			defer response.Body.Close()

			if response.StatusCode != tt.want.codeGet {
				t.Errorf("Expected status code %d, got %d", tt.want.codeGet, responseRecorder.Code)
			}

		})
	}
}
