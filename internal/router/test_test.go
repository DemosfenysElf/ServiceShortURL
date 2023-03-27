package router

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo"

	"ServiceShortURL/internal/shorturlservice"
)

func TestDB11111(t *testing.T) {

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
