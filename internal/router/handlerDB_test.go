package router

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
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
			for i, _ := range tt.generatorShort.Result {
				mock.ExpectExec("insert into ShortenerURL").
					WithArgs(tt.urllist[i], tt.generatorShort.ShortURL(), user.Name, user.Value, false).
					WillReturnResult(sqlmock.NewResult(1, 1))
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
