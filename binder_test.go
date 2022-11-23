package binder_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alikarimii/binder"
	"github.com/gorilla/mux"
)

type postReq struct {
	User     string `param:"user"`
	ID       int    `param:"id"`
	Email    string `form:"email"`
	Password string `form:"password"`
}
type getReq struct {
	User    string `param:"user"`
	ID      int    `param:"id"`
	IsAdmin bool   `query:"is_admin"`
	Test    int    `query:"test"`
}

func newHttpServer() *mux.Router {
	r := mux.NewRouter()

	r.Methods("POST").Path("/{user}/{id}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bi := binder.NewBinder()
		var mt postReq
		e := bi.Bind(&mt, r)
		if mt.ID != 1 || mt.Email == "" || mt.Password == "" { // test case
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(fmt.Errorf("cant get request body or param"))
			return
		}
		json.NewEncoder(w).Encode(e)
	})
	r.Methods("GET").Path("/{user}/{id}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bi := binder.NewBinder()
		var mt getReq
		e := bi.Bind(&mt, r)
		if mt.ID != 1 || !mt.IsAdmin || mt.Test != 2 { // test case
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(fmt.Errorf("cant get request query or param"))
			return
		}
		json.NewEncoder(w).Encode(e)
	})

	return r
}
func TestBind(t *testing.T) {

	srv := httptest.NewServer(newHttpServer())

	for _, testcase := range []struct {
		method, url, body string
		want              int
	}{
		{"POST", srv.URL + "/ali/1", `{"email": "ali@gmail.com", "password":"1234567"}`, http.StatusOK},
		{"POST", srv.URL + "/ali/1", "", http.StatusBadRequest}, // faild case
		{"GET", srv.URL + "/ali/1?is_admin=true&test=2", "", http.StatusOK},
		{"GET", srv.URL + "/ali/1", "", http.StatusBadRequest}, // faild case
	} {
		req, _ := http.NewRequest(testcase.method, testcase.url, strings.NewReader(testcase.body))
		req.Header.Set("Content-type", "application/json") // for body
		resp, _ := http.DefaultClient.Do(req)

		t.Run(fmt.Sprintf("start %s request", testcase.method), func(t *testing.T) {
			if testcase.want != resp.StatusCode {
				t.Errorf("%s %s %s: want %d have %d", testcase.method, testcase.url, testcase.body, testcase.want, resp.StatusCode)
			}
		})
	}

}

func TestWithMemory(t *testing.T) {

	t.Run("Create Binder with custom max memory", func(t *testing.T) {
		_ = binder.NewBinder(binder.WithCustomMemory(40 << 20)) // 40mb
		defer func() {
			if e := recover(); e != nil {
				t.Error("faild to create binder")
			}
		}()
	})

}

