package router

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo"

	"ServiceShortURL/internal/shorturlservice"
)

func TestDB111(t *testing.T) {

	type want struct {
		codeGet  int
		url      string
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
			errPing:        errors.New("Ping fail"),
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
