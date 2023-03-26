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

	"github.com/caarlos0/env"
	"github.com/labstack/echo"

	"ServiceShortURL/internal/shorturlservice"
)

func TestBatchGetDelete(t *testing.T) {

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
		baseurl map[string]string

		bodyURLJSON string
		urlResult   string
	}{
		{
			name: "TestApiShortenBatch1",
			want: want{
				codePost: 201,
				codeGet:  307,
				codeDel:  410,
				response: `{"status":"ok"}`,
			},
			baseurl: map[string]string{
				"baseurl1": "https://www.youtube.com/watch?v=UK7yzgVpnDA",
				"baseurl2": "https://metanit.com/",
				"baseurl3": "https://otion.site/context-19ec3dbb0cd4ec",
				"baseurl4": "https://stepik.org/lesson/359395",
			},

			bodyURLJSON: `[{"correlation_id":"baseurl1","original_url":"https://www.youtube.com/watch?v=UK7yzgVpnDA"},
								{"correlation_id":"baseurl2","original_url":"https://metanit.com/"},
								{"correlation_id":"baseurl3","original_url":"https://otion.site/context-19ec3dbb0cd4ec"},
								{"correlation_id":"baseurl4","original_url":"https://stepik.org/lesson/359395"}]`,
		},
		{
			name: "TestApiShortenBatch2",
			want: want{
				codePost: 500,
				response: `{"status":"ok"}`,
			},
			bodyURLJSON: "",
		},
		{
			name: "TestApiShortenBatch3",
			want: want{
				codePost: 500,
				response: `{"status":"ok"}`,
			},
			baseurl: map[string]string{
				"baseurl1": "https://www.2222.com/watch?v=UK7yzgVpnDA",
				"baseurl2": "https://4444.site/context-19ec3dbb0cd4ec",
			},
			bodyURLJSON: `[{"correlation_id":"baseurl1","original_url":https://www.2222.com/watch?v=UK7yzgVpnDA},
							{"correlation_id":baseurl2,"original_url":"https://4444.site/context-19ec3dbb0cd4ec"}]`,
		},
		{
			name: "TestApiShortenBatch4",
			want: want{
				codePost: 201,
				codeGet:  307,
				response: `{"status":"ok"}`,
			},
			baseurl: map[string]string{
				"baseurl1": "https://www.5555.com/watch?v=UK7yzgVpnDA",
				"baseurl2": "https://www.5555.com/watch?v=UK7yzgVpnDA",
			},
			bodyURLJSON: `[{"correlation_id":"baseurl1","original_url":"https://www.5555.com/watch?v=UK7yzgVpnDA"},
						   {"correlation_id":"baseurl2","original_url":"https://www.5555.com/watch?v=UK7yzgVpnDA"}]`,
		},
		{
			name: "TestApiShortenBatch5",
			want: want{
				codePost: 500,
				response: `{"status":"ok"}`,
			},
			bodyURLJSON: `[{"correlation_id":"baseurl1","original_url":"https://www.5555.com/watch?v=UK7yzgVpnD
						   {"correlation_id":"baseurl2","original_url":"https://www.5555.com/watch?v=UK7yzgVpnDA"]`,
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
				request := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", strings.NewReader(tt.bodyURLJSON))
				cookie := kyki()

				request.AddCookie(cookie)
				responseRecorder := httptest.NewRecorder()
				c := e.NewContext(request, responseRecorder)

				rout.StorageInterface = s

				rout.Serv = e
				rout.PostAPIShortenBatch(c)
				response := responseRecorder.Result()

				defer response.Body.Close()
				resBody, err := io.ReadAll(response.Body)
				if err != nil {
					t.Fatal(err)
				}

				if response.StatusCode != tt.want.codePost {
					t.Errorf("Expected status code %d, got %d", tt.want.codePost, responseRecorder.Code)
				}

				//////////////////////////////////////////////////////////////
				if response.StatusCode == http.StatusCreated {

					urlBatch := []shortURLApiShortenBatch{}
					err = json.Unmarshal(resBody, &urlBatch)
					if err != nil {
						t.Fatal(err)
					}
					str := `` + rout.Cfg.BaseURL + "/"
					s2 := `}`

					for i, batch := range urlBatch {

						resBodyShort := strings.Replace(batch.ShortURL, str, "", -1)
						resBodyShort = strings.Replace(resBodyShort, s2, "", -1)

						_, memory := rout.StorageInterface.(*shorturlservice.MemoryStorage)

						if (i == 2) && (!memory) {

							bodydelete := `["` + resBodyShort + `"]`
							request2 := httptest.NewRequest(http.MethodDelete, "/api/user/urls", strings.NewReader(bodydelete))
							request2.AddCookie(cookie)
							recorder2 := httptest.NewRecorder()
							ctx2 := e.NewContext(request2, recorder2)

							rout.DeleteAPIUserURLs(ctx2)
							time.Sleep(time.Second * 1)

						}

						request1 := httptest.NewRequest(http.MethodGet, "/"+resBodyShort, nil)
						recorder1 := httptest.NewRecorder()
						ctx1 := e.NewContext(request1, recorder1)

						rout.GetShortToURL(ctx1)
						response = recorder1.Result()
						defer response.Body.Close()
						resBodyGet, err := io.ReadAll(response.Body)
						if err != nil {
							t.Fatal(err)
						}

						if strings.HasPrefix(string(resBodyGet), rout.Cfg.BaseURL+"/") {
							t.Errorf("Expected body %s, got %s", tt.want.response, recorder1.Body.String())
						}

						if (i == 2) && (!memory) {
							if response.StatusCode != tt.want.codeDel {
								t.Errorf("Expected status code %d, got %d", tt.want.codeDel, recorder1.Code)
							}
							if response.Header.Get("Location") != "" {
								t.Errorf("Expected Location %s, got %s", "", response.Header.Get("Location"))
							}
						} else {
							if response.StatusCode != tt.want.codeGet {
								t.Errorf("Expected status code %d, got %d", tt.want.codeGet, recorder1.Code)
							}
							if response.Header.Get("Location") != tt.baseurl[batch.ID] {
								t.Errorf("Expected Location %s, got %s", tt.baseurl[batch.ID], response.Header.Get("Location"))
							}
						}
					}
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

func TestPostGet(t *testing.T) {
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
			name: "TestPostGet1",
			want: want{
				codePost: 201,
				codeGet:  307,
				response: `{"status":"ok"}`,
			},
			bodyURL: "https://www.1234.com/watch?v=UK7yzgVpnDA",
		},
		{
			name: "TestPostGet2",
			want: want{
				codePost: 201,
				codeGet:  307,
				response: `{"status":"ok"}`,
			},
			bodyURL: "https://www.4313.com/wyzgVpn",
		},
		{
			name: "TestPostGet3",
			want: want{
				codePost: 201,
				codeGet:  307,
				response: `{"status":"ok"}`,
			},
			bodyURL: "https://www.1234.com/watch?v=UK7yzgVpnDA",
		},
		{
			name: "TestPostGet4",
			want: want{
				codePost: 204,
				response: `{"status":"ok"}`,
			},
			bodyURL: "",
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
				responseRecorder := httptest.NewRecorder()
				c := e.NewContext(request, responseRecorder)

				rout.StorageInterface = s

				rout.Serv = e
				rout.PostURLToShort(c)
				response := responseRecorder.Result()

				defer response.Body.Close()
				resBody, err := io.ReadAll(response.Body)
				if err != nil {
					t.Fatal(err)
				}

				if response.StatusCode != tt.want.codePost {
					t.Errorf("Expected status code %d, got %d", tt.want.codePost, responseRecorder.Code)
				}

				//////////////////////////////////////////////////////////////
				if response.StatusCode == http.StatusCreated {

					str2 := `/`

					resBodyShort := strings.Replace(string(resBody), str2, "", -1)

					request1 := httptest.NewRequest(http.MethodGet, "/"+resBodyShort, nil)
					recorder1 := httptest.NewRecorder()
					ctx1 := e.NewContext(request1, recorder1)

					rout.GetShortToURL(ctx1)
					response = recorder1.Result()
					defer response.Body.Close()
					resBodyGet, err := io.ReadAll(response.Body)
					if err != nil {
						t.Fatal(err)
					}

					if strings.HasPrefix(string(resBodyGet), rout.Cfg.BaseURL+"/") {
						t.Errorf("Expected body %s, got %s", tt.want.response, recorder1.Body.String())
					}
					if response.StatusCode != tt.want.codeGet {
						t.Errorf("Expected status code %d, got %d", tt.want.codeGet, recorder1.Code)
					}
					if response.Header.Get("Location") != tt.bodyURL {
						t.Errorf("Expected Location %s, got %s", tt.bodyURL, response.Header.Get("Location"))
					}

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

			urlJSON: `{"url":""}`,
		},
		{
			name: "TestApiShorten3",
			want: want{
				codePost: 204,
				response: `{"status":"ok"}`,
			},
			urlJSON: ``,
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
				request := httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader(tt.urlJSON))
				rec := httptest.NewRecorder()
				c := e.NewContext(request, rec)
				rout.StorageInterface = s

				errConfig := env.Parse(&rout.Cfg)
				if errConfig != nil {
					t.Fatal(errConfig)
				}
				rout.Serv = e
				rout.PostAPIShorten(c)

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
				if res.StatusCode == http.StatusCreated {
					str := `{"result":"` + rout.Cfg.BaseURL + "/"
					str2 := `"}`

					resBodyShort := strings.Replace(string(resBody), str, "", -1)
					resBodyShort = strings.Replace(resBodyShort, str2, "", -1)

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
			}
			teststor(&shorturlservice.FileStorage{
				FilePath:    rout.Cfg.Storage,
				RandomShort: &shorturlservice.RandomGenerator{},
			})
			teststor(shorturlservice.InitMem())

		})
	}
}

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
				FilePath:    rout.Cfg.Storage,
				RandomShort: &shorturlservice.RandomGenerator{},
			})

		})
	}
}

func Test_PostGet(t *testing.T) {

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
			name: "Test_PostGet_1",
			want: want{
				codePost: 201,
				codeGet:  307,
				response: `{"status":"ok"}`,
			},
			url: ("https://www.youtube.com/watch?v=UK7yzgVpnDA"),
		}, {
			name: "Test_PostGet_2",
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
			//os.Truncate(testStorageURL, 0)
			//os.Truncate(testStorageUsers, 0)

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
				defer resGetu2.Body.Close()
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
				defer resDel.Body.Close()
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
				defer resGet.Body.Close()
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
				defer resDelu2.Body.Close()
				if resDelu2.StatusCode != tt.want.codeDel2 {
					t.Errorf("Expected status code %d, got %d", tt.want.codeDel2, responseRecorderDelu2.Code)
				}

			}

			teststor(&shorturlservice.FileStorage{
				FilePath:    rout.Cfg.Storage,
				RandomShort: &shorturlservice.RandomGenerator{},
			})

		})
	}
}
