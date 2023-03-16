package router

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo"

	"ServiceShortURL/internal/shorturlservice"
)

func TestDeleteGetu1u2(t *testing.T) {

	type want struct {
		codePost int
		codeGet1 int
		codeGet2 int
		codeDel1 int
		codeDel2 int
		response string
	}
	tests := []struct {
		name        string
		want        want
		bodyURLJSON string
	}{
		{
			name: "TestDeleteGetu1u2",
			want: want{
				codePost: 201,
				codeGet1: 410,
				codeGet2: 307,
				codeDel1: 202,
				codeDel2: 202,
				response: `{"status":"ok"}`,
			},

			bodyURLJSON: `[{"correlation_id":"baseurl1","original_url":"https://www.youtube.com/watch?v=UK7yzgVpnDA"},
								{"correlation_id":"baseurl2","original_url":"https://metanit.com/"},
								{"correlation_id":"baseurl3","original_url":"https://otion.site/context-19ec3dbb0cd4ec"},
								{"correlation_id":"baseurl4","original_url":"https://stepik.org/lesson/359395"}]`,
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
				rout.StorageInterface = s
				rout.Serv = e
				// записываем ссылки от user1 через batch
				var user1 = kyki()
				request := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", strings.NewReader(tt.bodyURLJSON))
				request.AddCookie(user1)
				responseRecorder := httptest.NewRecorder()
				ctx := e.NewContext(request, responseRecorder)
				rout.PostAPIShortenBatch(ctx)

				// обрабатываем результат
				res := responseRecorder.Result()
				defer res.Body.Close()
				resBody, err := io.ReadAll(res.Body)
				if err != nil {
					t.Fatal(err)
				}
				if res.StatusCode != tt.want.codePost {
					t.Errorf("Expected status code %d, got %d", tt.want.codePost, http.StatusCreated)
				}
				urlBatch := []shortURLApiShortenBatch{}
				err = json.Unmarshal(resBody, &urlBatch)
				if err != nil {
					t.Fatal(err)
				}
				str := `` + rout.Cfg.BaseURL + "/"
				str2 := `}`
				resBodyShort1 := strings.Replace(urlBatch[1].ShortURL, str, "", -1)
				resBodyShort1 = strings.Replace(resBodyShort1, str2, "", -1)
				resBodyShort2 := strings.Replace(urlBatch[2].ShortURL, str, "", -1)
				resBodyShort2 = strings.Replace(resBodyShort2, str2, "", -1)
				time.Sleep(1 * time.Second)
				// get [2] третью ссылку ВТОРЫМ пользователем
				var user2 = kyki()
				requestGetu2 := httptest.NewRequest(http.MethodGet, "/"+resBodyShort2, nil)
				requestGetu2.AddCookie(user2)
				responseRecorderGetu2 := httptest.NewRecorder()
				ctxGetu2 := e.NewContext(requestGetu2, responseRecorderGetu2)
				rout.GetShortToURL(ctxGetu2)
				resGetu2 := responseRecorderGetu2.Result()
				if resGetu2.StatusCode != tt.want.codeGet2 {
					t.Errorf("Expected status code %d, got %d", tt.want.codeGet2, responseRecorderGetu2.Code)
				}

				// delete [1] вторую ссылку первым пользователем
				resBodyShortDel1 := `["` + resBodyShort1 + `"]`
				requestDel := httptest.NewRequest(http.MethodDelete, "/api/user/urls", strings.NewReader(resBodyShortDel1))
				requestDel.AddCookie(user1)
				responseRecorderDel := httptest.NewRecorder()
				ctxDel := e.NewContext(requestDel, responseRecorderDel)
				rout.DeleteAPIUserURLs(ctxDel)
				resDel := responseRecorderDel.Result()
				if resDel.StatusCode != tt.want.codeDel1 {
					t.Errorf("Expected status code %d, got %d", tt.want.codeDel1, responseRecorderDel.Code)
				}
				time.Sleep(1 * time.Second)
				// get [1] вторую ссылку первым пользователем
				requestGet := httptest.NewRequest(http.MethodGet, "/"+resBodyShort1, nil)
				requestGet.AddCookie(user1)
				responseRecorderGet := httptest.NewRecorder()
				ctxGetu1 := e.NewContext(requestGet, responseRecorderGet)
				rout.GetShortToURL(ctxGetu1)
				resGet := responseRecorderGet.Result()
				if resGet.StatusCode != tt.want.codeGet1 {
					t.Errorf("Expected status code %d, got %d", tt.want.codeGet1, responseRecorderGet.Code)
				}

				// delete [2] третью ссылку ВТОРЫМ пользователем
				resBodyShortDel2 := `["` + resBodyShort2 + `"]`
				requestDelu2 := httptest.NewRequest(http.MethodDelete, "/api/user/urls", strings.NewReader(resBodyShortDel2))
				requestDelu2.AddCookie(user2)
				responseRecorderDelu2 := httptest.NewRecorder()
				ctxDelu2 := e.NewContext(requestDelu2, responseRecorderDelu2)
				rout.DeleteAPIUserURLs(ctxDelu2)
				resDelu2 := responseRecorderDelu2.Result()
				if resDelu2.StatusCode != tt.want.codeDel2 {
					t.Errorf("Expected status code %d, got %d", tt.want.codeDel2, responseRecorderDelu2.Code)
				}

			}

			teststor(&shorturlservice.FileStorage{
				FilePath: rout.Cfg.Storage,
			})

		})
	}
}
