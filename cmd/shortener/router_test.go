package main

import (
	"ServiceShortURL/internal/router"
	"github.com/caarlos0/env"
	"github.com/labstack/echo"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_router(t *testing.T) {

	type want struct {
		codePost    int
		codeGet     int
		response    string
		contentType string
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
				codePost: 201,
				codeGet:  400,
				response: `{"status":"ok"}`,
			},
			url: (""),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.url))
			rec := httptest.NewRecorder()
			c := e.NewContext(request, rec)

			rout := router.Server{}
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

		})
	}
}
