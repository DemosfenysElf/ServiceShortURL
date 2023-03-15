package router

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/labstack/echo"

	"ServiceShortURL/internal/shorturlservice"
)

func TestGetApiUserURL(t *testing.T) {

	var bodyURLJSON2 = `[{"correlation_id":"baseurl11","original_url":"https://www.00outube.com/watch?v=UK7yzgVpnDA"},
						{"correlation_id":"baseurl12","original_url":"https://99metanit.com/"},
						{"correlation_id":"baseurl13","original_url":"https://88otion.site/context-19ec3dbb0cd4ec"},
						{"correlation_id":"baseurl14","original_url":"https://77stepik.org/lesson/359395"}]`
	type want struct {
		codePost int
		codeGet  int
		codeDel  int
		response string
	}
	tests := []struct {
		name    string
		want    want
		wantErr bool
		baseurl map[int]string

		bodyURLJSON string
		urlResult   string
	}{
		{
			name: "TestGetApiUserURL1",
			want: want{
				codePost: 201,
				codeGet:  307,
				codeDel:  410,
				response: `{"status":"ok"}`,
			},
			baseurl: map[int]string{
				1: "https://www.youtube.com/watch?v=UK7yzgVpnDA",
				2: "https://metanit.com/",
				3: "https://otion.site/context-19ec3dbb0cd4ec",
				4: "https://stepik.org/lesson/359395",
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
				// записываем ссылки от user1 через batch
				var user1 = kyki()
				request := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", strings.NewReader(tt.bodyURLJSON))
				request.AddCookie(user1)
				responseRecorder := httptest.NewRecorder()
				ctx := e.NewContext(request, responseRecorder)
				rout.StorageInterface = s
				rout.Serv = e
				rout.PostAPIShortenBatch(ctx)

				// записываем ссылки от user2 через batch
				var user2 = kyki()
				request2 := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", strings.NewReader(bodyURLJSON2))
				request2.AddCookie(user2)
				responseRecorder2 := httptest.NewRecorder()
				ctx2 := e.NewContext(request2, responseRecorder2)
				rout.StorageInterface = s
				rout.Serv = e
				rout.PostAPIShortenBatch(ctx2)

				// читаем записанные ссылки полученные от пользователя user1
				request3 := httptest.NewRequest(http.MethodPost, "/api/user/urls", nil)
				request3.AddCookie(user1)
				responseRecorder = httptest.NewRecorder()
				ctx = e.NewContext(request3, responseRecorder)
				rout.GetAPIUserURL(ctx)
				response := responseRecorder.Result()
				defer response.Body.Close()
				resBody3, err := io.ReadAll(response.Body)
				if err != nil {
					t.Fatal(err)
				}

				urlBatch := []userURLstruct{}
				err = json.Unmarshal(resBody3, &urlBatch)
				if err != nil {
					t.Fatal(err)
				}
				i := 0
				for _, i2 := range urlBatch {
					switch i2.OriginalURL {
					case tt.baseurl[1]:
						i += 1
					case tt.baseurl[2]:
						i += 10
					case tt.baseurl[3]:
						i += 100
					case tt.baseurl[4]:
						i += 1000
					default:
						t.Errorf("Неподходящяя ссылка %s", i2.OriginalURL)
					}
				}
				if i != 1111 {
					t.Errorf("Несоответстветствующее количество ссылок")
				}

				// получение данных user3 (без данных)
				var user3 = kyki()
				request4 := httptest.NewRequest(http.MethodPost, "/api/user/urls", nil)
				request4.AddCookie(user3)
				responseRecorder = httptest.NewRecorder()
				ctx = e.NewContext(request4, responseRecorder)
				rout.GetAPIUserURL(ctx)
				response = responseRecorder.Result()
				defer response.Body.Close()
				resBody4, err := io.ReadAll(response.Body)
				if err != nil {
					t.Fatal(err)
				}
				if len(resBody4) != 0 {
					t.Errorf("Получены данные, их быть недолжно")
				}
			}

			teststor(&shorturlservice.FileStorage{
				FilePath: rout.Cfg.Storage,
			})

		})
	}
}
