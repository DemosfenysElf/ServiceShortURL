package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_handler(t *testing.T) {

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
			name: "test_1",
			want: want{
				codePost: 201,
				codeGet:  307,
				response: `{"status":"ok"}`,
			},
			url: ("https://www.youtube.com/watch?v=UK7yzgVpnDA"),
		}, {
			name: "test_2",
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
			w := httptest.NewRecorder()
			h := http.HandlerFunc(handler)

			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.url))
			h.ServeHTTP(w, request)
			res := w.Result()

			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}

			if res.StatusCode != tt.want.codePost {
				t.Errorf("Expected status code %d, got %d", tt.want.codePost, w.Code)
			}

			resBodyShort := strings.Replace(string(resBody), "http://localhost:8080/", "", -1)
			//resBodyIoReader := strings.NewReader(resBodyShort)
			request1 := httptest.NewRequest(http.MethodGet, "/"+resBodyShort, nil)
			w1 := httptest.NewRecorder()
			h.ServeHTTP(w1, request1)
			res = w1.Result()

			defer res.Body.Close()
			resBody, err = io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}

			if strings.HasPrefix(string(resBody), "http://localhost:8080/") {
				t.Errorf("Expected body %s, got %s", tt.want.response, w1.Body.String())
			}

			if res.StatusCode != tt.want.codeGet {
				t.Errorf("Expected status code %d, got %d", tt.want.codeGet, w1.Code)
			}

			if res.Header.Get("Location") != tt.url {
				t.Errorf("Expected Location %s, got %s", tt.url, res.Header.Get("Location"))
			}

		})
	}
}
