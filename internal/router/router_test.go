package router

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/caarlos0/env"
	"github.com/labstack/echo"

	"ServiceShortURL/internal/shorturlservice"
)

func Test_router(t *testing.T) {

	type want struct {
		codePost int
		codeGet  int
		response string
	}

	tests := []struct {
		name string
		want want
		url  string
	}{
		{
			name: "Test_router_1",
			want: want{
				codePost: 201,
				codeGet:  307,
				response: `{"status":"ok"}`,
			},
			url: ("https://www.youtube.com/watch?v=UK7yzgVpnDA"),
		}, {
			name: "Test_router_2",
			want: want{
				codePost: 204,
				response: `{"status":"ok"}`,
			},
			url: (""),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Truncate(testStorageURL, 0)
			os.Truncate(testStorageUsers, 0)
			e := echo.New()
			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.url))
			rec := httptest.NewRecorder()
			c := e.NewContext(request, rec)

			rout := InitServer()
			rout.StorageInterface = shorturlservice.InitMem()

			errConfig := env.Parse(&rout.Cfg)
			if errConfig != nil {
				t.Fatal(errConfig)
			}
			rout.Serv = e
			rout.PostURLToShort(c)

			res := rec.Result()

			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}

			if res.StatusCode != tt.want.codePost {
				t.Errorf("Expected status code %d, got %d", tt.want.codePost, rec.Code)
			}
			if res.StatusCode == http.StatusCreated {
				//////////////////////////////////////////////////////////////

				resBodyShort := strings.Replace(string(resBody), rout.Cfg.BaseURL+"/", "", -1)
				request1 := httptest.NewRequest(http.MethodGet, "/"+resBodyShort, nil)

				rec1 := httptest.NewRecorder()
				c1 := e.NewContext(request1, rec1)

				rout.GetShortToURL(c1)
				res = rec1.Result()

				defer res.Body.Close()
				resBody, err = io.ReadAll(res.Body)
				if err != nil {
					t.Fatal(err)
				}

				if strings.HasPrefix(string(resBody), rout.Cfg.BaseURL+"/") {
					t.Errorf("Expected body %s, got %s", tt.want.response, rec1.Body.String())
				}

				if res.StatusCode != tt.want.codeGet {
					t.Errorf("Expected status code %d, got %d", tt.want.codeGet, rec1.Code)
				}

				if res.Header.Get("Location") != tt.url {
					t.Errorf("Expected Location %s, got %s", tt.url, res.Header.Get("Location"))
				}
			}
		})
	}
}