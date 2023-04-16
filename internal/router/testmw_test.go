package router

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo"

	"ServiceShortURL/internal/shorturlservice"
)

// будут тесты авторизации как пойму как их делать

func TestAutPost1(t *testing.T) {
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
			name: "TestAut",
			want: want{
				codePost: 201,
				response: `{"status":"ok"}`,
			},
			bodyURL: "https://www.523.com/watch?v=UK7yzgVpnDA",
		},
		{
			name: "TestAut",
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
				e := echo.New()
				request := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(tt.bodyURL)))
				cookie := kyki()
				request.AddCookie(cookie)
				responseRecorder := httptest.NewRecorder()
				//c := e.NewContext(request, responseRecorder)
				rout.StorageInterface = s

				////////
				e.Use(rout.MWAuthentication)
				e.POST("/", rout.PostURLToShort)
				e.ServeHTTP(responseRecorder, request)
				/////////

				//rout.Serv = e
				//rout.PostURLToShort(c)
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
			//teststor(shorturlservice.InitMem())
		})
	}
}

//func TestName(t *testing.T) {
//	type want struct {
//		codePost int
//		codeGet  int
//		response string
//	}
//	tests := []struct {
//		name    string
//		want    want
//		wantErr bool
//		bodyURL string
//	}{
//		{
//			name: "TestPostGetGZ",
//			want: want{
//				codePost: 201,
//				codeGet:  307,
//				response: `{"status":"ok"}`,
//			},
//			bodyURL: "https://www.1234.com/watch?v=UK7yzgVpnDA",
//		},
//	}
//	for _ = range tests {
//		// Инициализация
//		rout := InitServer()
//		//rootUrl := func(c echo.Context) error {
//		//	return c.HTML(http.StatusOK, "body output")
//		//}
//		e := echo.New()
//		e.Use(rout.mwGzipHandle)
//		e.GET("/", rout.PostURLToShort)
//
//		// Запрос
//		req := httptest.NewRequest(echo.GET, "/", nil)
//		rec := httptest.NewRecorder()
//		// Using the ServerHTTP on echo will trigger the router and middleware
//		e.ServeHTTP(rec, req)
//
//		// Проверка работы
//
//	}
//}
