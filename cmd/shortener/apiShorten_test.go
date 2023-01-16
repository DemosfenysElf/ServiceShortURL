package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"ServiceShortURL/internal/router"
	"ServiceShortURL/internal/shorturlservice"

	"github.com/caarlos0/env"
	"github.com/labstack/echo"
)

func TestApiShorten(t *testing.T) {

	type want struct {
		codePost int
		codeGet  int
		response string
	}
	tests := []struct {
		name string

		want    want
		wantErr bool
		baseurl string
		urlJSON string
		urlRes  string
	}{
		{
			name: "TestApiShorten1",
			want: want{
				codePost: 201,
				codeGet:  307,
				response: `{"status":"ok"}`,
			},
			baseurl: "https://www.youtube.com/watch?v=UK7yzgVpnDA",
			urlJSON: `{"url":"https://www.youtube.com/watch?v=UK7yzgVpnDA"}`,
			urlRes:  `{"result":"https://www.youtube.com/watch?v=UK7yzgVpnDA"}`,
		},
		{
			name: "TestApiShorten2",
			want: want{
				codePost: 400,
				response: `{"status":"ok"}`,
			},
			baseurl: "",
			urlJSON: `{"url":""}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			e := echo.New()
			request := httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader(tt.urlJSON))
			rec := httptest.NewRecorder()
			c := e.NewContext(request, rec)
			rout := router.URLServer{}
			rout.StorageInterface = shorturlservice.InitMem()

			errConfig := env.Parse(&rout.Cfg)
			if errConfig != nil {
				t.Fatal(errConfig)
			}
			rout.Serv = e
			rout.APIShorten(c)

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
			if res.StatusCode == 201 {
				s := `{"result":"` + rout.Cfg.BaseURL + "/"
				s2 := `"}`

				resBodyShort := strings.Replace(string(resBody), s, "", -1)
				resBodyShort = strings.Replace(resBodyShort, s2, "", -1)

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

				if res.Header.Get("Location") != tt.baseurl {
					t.Errorf("Expected Location %s, got %s", tt.baseurl, res.Header.Get("Location"))
				}
			}
		})
	}
}
