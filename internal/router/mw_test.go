package router

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo"

	"ServiceShortURL/internal/shorturlservice"
)

func TestMWPostGZIP(t *testing.T) {
	type want struct {
		codePost int
		response string
	}
	tests := []struct {
		name    string
		want    want
		wantErr bool
		bodyURL string
	}{
		{
			name: "TestMWPostGZIP",
			want: want{
				codePost: 201,
				response: `{"status":"ok"}`,
			},
			bodyURL: "https://www.1234.com/watch?v=UK7yzgVpnDA",
		},
		{
			name: "TestPostMWGZIP2",
			want: want{
				codePost: 204,
				response: `{"status":"ok"}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Truncate(testStorageURL, 0)
			os.Truncate(testStorageUsers, 0)
			rout := InitServer()
			rout.Cfg.Storage = testStorageURL

			teststor := func(s shorturlservice.StorageInterface) {
				compressBody, _ := shorturlservice.ServiceCompress([]byte(tt.bodyURL))
				e := echo.New()
				request := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(compressBody))
				cookie := kyki()
				request.AddCookie(cookie)
				request.Header.Add("Accept-Encoding", "gzip")
				request.Header.Add("Content-Encoding", "gzip")
				responseRecorder := httptest.NewRecorder()
				rout.StorageInterface = s

				e.Use(rout.mwGzipHandle)
				e.POST("/", rout.PostURLToShort)
				e.ServeHTTP(responseRecorder, request)

				response := responseRecorder.Result()
				defer response.Body.Close()
				if response.StatusCode != tt.want.codePost {
					t.Errorf("Expected status code %d, got %d", tt.want.codePost, responseRecorder.Code)
				}

			}
			teststor(&shorturlservice.FileStorage{
				FilePath:    rout.Cfg.Storage,
				RandomShort: &shorturlservice.RandomGenerator{},
			})
			teststor(shorturlservice.InitMem())
		})
	}
}

func TestAutPost(t *testing.T) {
	type want struct {
		codePost int
		response string
	}
	tests := []struct {
		name      string
		want      want
		wantErr   bool
		bodyURL   string
		badCookie bool
	}{
		{
			name: "TestAut1",
			want: want{
				codePost: 201,
				response: `{"status":"ok"}`,
			},
			bodyURL:   "https://www.523.com/watch?v=UK7yzgVpnDA",
			badCookie: false,
		},
		{
			name: "TestAut2",
			want: want{
				codePost: 201,
				response: `{"status":"ok"}`,
			},
			bodyURL:   "https://www.523.com/watch?v=UK7yzgVpnDA",
			badCookie: true,
		},
		{
			name: "TestAut3",
			want: want{
				codePost: 204,
				response: `{"status":"ok"}`,
			},
			badCookie: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Truncate(testStorageURL, 0)
			rout := InitServer()
			rout.Cfg.Storage = testStorageURL

			teststor := func(s shorturlservice.StorageInterface) {
				e := echo.New()
				request := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(tt.bodyURL)))
				cookie := kyki()
				if tt.badCookie {
					cookie.Name = cookie.Name[1:]
				}
				request.AddCookie(cookie)
				responseRecorder := httptest.NewRecorder()
				rout.StorageInterface = s

				e.Use(rout.MWAuthentication)
				e.POST("/", rout.PostURLToShort)
				e.ServeHTTP(responseRecorder, request)

				response := responseRecorder.Result()
				defer response.Body.Close()
				if response.StatusCode != tt.want.codePost {
					t.Errorf("Expected status code %d, got %d", tt.want.codePost, responseRecorder.Code)
				}
			}
			teststor(&shorturlservice.FileStorage{
				FilePath:    rout.Cfg.Storage,
				RandomShort: &shorturlservice.RandomGenerator{},
			})
			teststor(shorturlservice.InitMem())
		})
	}
}

func TestIPStats(t *testing.T) {
	tests := []struct {
		name           string
		codeGet        int
		generatorShort *shorturlservice.TestGenerator
		generatorUser  *shorturlservice.TestGeneratorUser
		want           string
		ip             string
	}{
		{
			name:           "TestIPStats1",
			codeGet:        200,
			generatorShort: &shorturlservice.TestGenerator{Result: []string{"Short1"}, Index: 0},
			generatorUser:  &shorturlservice.TestGeneratorUser{Result: []byte("UserUserUserUser")},
			want:           `{"urls":123,"users":1}`,
			ip:             "192.0.2.1/24",
		},
		{
			name:           "TestIPStats2",
			codeGet:        403,
			generatorShort: &shorturlservice.TestGenerator{Result: []string{"Short1"}, Index: 0},
			generatorUser:  &shorturlservice.TestGeneratorUser{Result: []byte("UserUserUserUser")},
			ip:             "1927.2.1/24",
		},
		{
			name:           "TestIPStats3",
			codeGet:        403,
			generatorShort: &shorturlservice.TestGenerator{Result: []string{"Short1"}, Index: 0},
			generatorUser:  &shorturlservice.TestGeneratorUser{Result: []byte("UserUserUserUser")},
			ip:             "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Truncate(testStorageURL, 0)
			os.Truncate(testStorageUsers, 0)
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()
			rout := InitTestServer()
			rout.Cfg.Storage = testStorageURL
			rout.Cfg.TrustedSubnet = tt.ip
			mockDB := &shorturlservice.Database{RandomShort: tt.generatorShort}

			mockDB.SetConnection(db)
			rout.DB = mockDB
			rout.StorageInterface = mockDB
			rout.GeneratorUsers = tt.generatorUser
			user := rout.kykiDB()

			rs := sqlmock.NewRows([]string{"count"}).AddRow(123)
			mock.ExpectQuery(regexp.QuoteMeta("select count(url) from ShortenerURL")).
				WillReturnRows(rs)
			tt.generatorShort.Index = 0

			e := echo.New()
			e.Use(rout.MWCheakerIP)
			e.GET("/api/internal/stats", rout.GetAPIInternalStats)
			request := httptest.NewRequest(http.MethodGet, "/api/internal/stats", nil)
			request.Header.Add("X-Real-IP", "192.0.2.1")
			request.AddCookie(user)
			responseRecorder := httptest.NewRecorder()
			e.ServeHTTP(responseRecorder, request)

			rout.Serv = e

			response := responseRecorder.Result()

			response.Header.Get("Location")
			defer response.Body.Close()

			if response.StatusCode != tt.codeGet {
				t.Errorf("Expected status code %d, got %d", tt.codeGet, responseRecorder.Code)
			}
			body, _ := io.ReadAll(response.Body)
			if string(body) != tt.want {
				t.Errorf("Expected body %s, got %s", tt.want, string(body))
			}
		})
	}
}
