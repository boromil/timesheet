package slashdb

import (
	"net/http"
	"net/http/httptest"
)

type testDoer struct {
	respBody       string
	respStatusCode int
}

func (d *testDoer) Do(req *http.Request) (*http.Response, error) {
	w := httptest.NewRecorder()
	w.Write([]byte(d.respBody))
	w.Code = d.respStatusCode
	w.HeaderMap.Set("Content-Type", "application/json")
	return w.Result(), nil
}
