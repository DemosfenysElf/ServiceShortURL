package test

// func TestApiShortenBatch(t *testing.T) {
//
//		type want struct {
//			codePost int
//			codeGet  int
//			codeDel  int
//			response string
//		}
//		tests := []struct {
//			name    string
//			want    want
//			wantErr bool
//			baseurl map[string]string
//
//			bodyURLJSON string
//			urlResult   string
//		}{
//			{
//				name: "TestApiShortenBatch1",
//				want: want{
//					codePost: 201,
//					codeGet:  307,
//					codeDel:  410,
//					response: `{"status":"ok"}`,
//				},
//				baseurl: map[string]string{
//					"baseurl1": "https://www.youtube.com/watch?v=UK7yzgVpnDA",
//					"baseurl2": "https://metanit.com/",
//					"baseurl3": "https://otion.site/context-19ec3dbb0cd4ec",
//					"baseurl4": "https://stepik.org/lesson/359395",
//				},
//
//				bodyURLJSON: `[{"correlation_id":"baseurl1","original_url":"https://www.youtube.com/watch?v=UK7yzgVpnDA"},
//								{"correlation_id":"baseurl2","original_url":"https://metanit.com/"},
//								{"correlation_id":"baseurl3","original_url":"https://otion.site/context-19ec3dbb0cd4ec"},
//								{"correlation_id":"baseurl4","original_url":"https://stepik.org/lesson/359395"}]`,
//			},
//		}
//		for _, tt := range tests {
//			t.Run(tt.name, func(t *testing.T) {
//
//				e := echo.New()
//				request := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", strings.NewReader(tt.bodyURLJSON))
//				cookie := kyki()
//
//				request.AddCookie(cookie)
//				responseRecorder := httptest.NewRecorder()
//				c := e.NewContext(request, responseRecorder)
//				rout := router.InitServer()
//				rout.InitRouter()
//
//				DB := &shorturlservice.Database{}
//				if errConnect := DB.Connect(rout.Cfg.ConnectDB); errConnect != nil {
//					t.Errorf("DB.Connect fail")
//				}
//				rout.StorageInterface = DB
//
//				rout.Serv = e
//				rout.PostAPIShortenBatch(c)
//				response := responseRecorder.Result()
//
//				defer response.Body.Close()
//				resBody, err := io.ReadAll(response.Body)
//				if err != nil {
//					t.Fatal(err)
//				}
//
//				if response.StatusCode != tt.want.codePost {
//					t.Errorf("Expected status code %d, got %d", tt.want.codePost, responseRecorder.Code)
//				}
//
//				//////////////////////////////////////////////////////////////
//				if response.StatusCode == 201 {
//
//					urlBatch := []shortURLApiShortenBatch{}
//					err = json.Unmarshal(resBody, &urlBatch)
//					if err != nil {
//						t.Fatal(err)
//					}
//					s := `` + rout.Cfg.BaseURL + "/"
//					s2 := `}`
//
//					for i, batch := range urlBatch {
//
//						resBodyShort := strings.Replace(batch.ShortURL, s, "", -1)
//						resBodyShort = strings.Replace(resBodyShort, s2, "", -1)
//
//						if i == 2 {
//							bodydelete := `["` + resBodyShort + `"]`
//							request2 := httptest.NewRequest(http.MethodDelete, "/api/user/urls", strings.NewReader(bodydelete))
//							request2.AddCookie(cookie)
//							recorder2 := httptest.NewRecorder()
//							ctx2 := e.NewContext(request2, recorder2)
//
//							rout.DeleteAPIUserURL(ctx2)
//
//						}
//
//						request1 := httptest.NewRequest(http.MethodGet, "/"+resBodyShort, nil)
//						recorder1 := httptest.NewRecorder()
//						ctx1 := e.NewContext(request1, recorder1)
//
//						rout.GetShortToURL(ctx1)
//						response = recorder1.Result()
//						defer response.Body.Close()
//						resBody, err = io.ReadAll(response.Body)
//						if err != nil {
//							t.Fatal(err)
//						}
//
//						if strings.HasPrefix(string(resBody), rout.Cfg.BaseURL+"/") {
//							t.Errorf("Expected body %s, got %s", tt.want.response, recorder1.Body.String())
//						}
//
//						if i == 2 {
//							if response.StatusCode != tt.want.codeDel {
//								t.Errorf("Expected status code %d, got %d", tt.want.codeDel, recorder1.Code)
//							}
//							if response.Header.Get("Location") != "" {
//								t.Errorf("Expected Location %s, got %s", "", response.Header.Get("Location"))
//							}
//						} else {
//							if response.StatusCode != tt.want.codeGet {
//								t.Errorf("Expected status code %d, got %d", tt.want.codeGet, recorder1.Code)
//							}
//							if response.Header.Get("Location") != tt.baseurl[batch.ID] {
//								t.Errorf("Expected Location %s, got %s", tt.baseurl[batch.ID], response.Header.Get("Location"))
//							}
//						}
//
//					}
//				}
//			})
//		}
//	}
