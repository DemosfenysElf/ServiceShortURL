package router

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/labstack/echo"

	"ServiceShortURL/internal/shorturlservice"
)

func TestPostGZIPOut(t *testing.T) {
	type want struct {
		codePost int
		codeGet  int
		response string
	}
	tests := []struct {
		name    string
		want    want
		wantErr bool
		bodyURL string
	}{
		{
			name: "TestPostGetGZ",
			want: want{
				codePost: 201,
				codeGet:  307,
				response: `{"status":"ok"}`,
			},
			bodyURL: "https://www.1234.com/watch?v=UK7yzgVpnDA",
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
				request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.bodyURL))
				cookie := kyki()
				request.AddCookie(cookie)
				request.Header.Add("Accept-Encoding", "gzip")
				responseRecorder := httptest.NewRecorder()
				c := e.NewContext(request, responseRecorder)

				rout.StorageInterface = s

				rout.Serv = e
				rout.PostURLToShort(c)
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
