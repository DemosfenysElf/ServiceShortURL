package router

// будут тесты авторизации как пойму как их делать

//func TestName(t *testing.T) {
//	type want struct {
//		codePost int
//		codeGet  int
//		response string
//	}
//	tests := []struct {
//		name    string
//		want    want
//		wantErr bool
//		bodyURL string
//	}{
//		{
//			name: "TestPostGetGZ",
//			want: want{
//				codePost: 201,
//				codeGet:  307,
//				response: `{"status":"ok"}`,
//			},
//			bodyURL: "https://www.1234.com/watch?v=UK7yzgVpnDA",
//		},
//	}
//	for _ = range tests {
//		// Инициализация
//		rout := InitServer()
//		//rootUrl := func(c echo.Context) error {
//		//	return c.HTML(http.StatusOK, "body output")
//		//}
//		e := echo.New()
//		e.Use(rout.mwGzipHandle)
//		e.GET("/", rout.PostURLToShort)
//
//		// Запрос
//		req := httptest.NewRequest(echo.GET, "/", nil)
//		rec := httptest.NewRecorder()
//		// Using the ServerHTTP on echo will trigger the router and middleware
//		e.ServeHTTP(rec, req)
//
//		// Проверка работы
//
//	}
//}
