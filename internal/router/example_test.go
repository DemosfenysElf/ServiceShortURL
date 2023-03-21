package router

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/labstack/echo"

	"ServiceShortURL/internal/shorturlservice"
)

var fileStorage = "../test/shortsURl.log"

func ExampleserverShortener_PostAPIShortenBatch() {

	baseurl := map[string]string{
		"baseurl1": "https://www.youtube.com/watch?v=UK7yzgVpnDA",
		"baseurl2": "https://metanit.com/",
		"baseurl3": "https://otion.site/context-19ec3dbb0cd4ec",
		"baseurl4": "https://stepik.org/lesson/359395",
	}
	bodyURLJSON := `[{"correlation_id":"baseurl1","original_url":"https://www.youtube.com/watch?v=UK7yzgVpnDA"},
					{"correlation_id":"baseurl2","original_url":"https://metanit.com/"},
					{"correlation_id":"baseurl3","original_url":"https://otion.site/context-19ec3dbb0cd4ec"},
					{"correlation_id":"baseurl4","original_url":"https://stepik.org/lesson/359395"}]`
	//

	e := echo.New()
	request := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", strings.NewReader(bodyURLJSON))
	cookie := kyki()

	request.AddCookie(cookie)
	responseRecorder := httptest.NewRecorder()
	c := e.NewContext(request, responseRecorder)
	rout := InitServer()
	rout.StorageInterface = &shorturlservice.FileStorage{
		FilePath: fileStorage,
	}

	rout.Serv = e
	rout.PostAPIShortenBatch(c)
	response := responseRecorder.Result()

	defer response.Body.Close()
	resBody, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}

	if response.StatusCode == http.StatusCreated {

		urlBatch := []shortURLApiShortenBatch{}
		err = json.Unmarshal(resBody, &urlBatch)
		if err != nil {
			fmt.Println(err)
		}

		for _, b := range urlBatch {
			fmt.Println("Original URL " + baseurl[b.ID])
		}

		// Output:
		// ==>> APIShortenBatch
		// Original URL https://www.youtube.com/watch?v=UK7yzgVpnDA
		// Original URL https://metanit.com/
		// Original URL https://otion.site/context-19ec3dbb0cd4ec
		// Original URL https://stepik.org/lesson/359395

	}
}
