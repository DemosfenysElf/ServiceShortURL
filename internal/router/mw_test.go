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
