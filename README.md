## Binder
binder is a form/param/query/header converter of http request to struct.this package is custom version of bind in [echo](https://echo.labstack.com/) framework for [mux](https://github.com/gorilla/mux) router.

### install

`go get github.com/alikarimii/binder`

### How it's work

```
package main

import (
  "github.com/alikarimii/binder"
  "github.com/gorilla/mux"
)

type body struct {
    Name string `form:"name"`
    Age int `form:"age"`
}

type param struct {
    User string `param:"user"`
    ID   int `param:"id"`
}
type query struct {
    IsAdmin bool `query:"is_admin"`
    Test int `query:"test"` 
}

func main() {
    b := binder.NewBinder()
    r := mux.NewRouter()
    r.Methods("POST").Path("/api").HandlerFunc(func(w http.ResponseWriter, r *http.Request){
       var mystruct body
       err := b.Bind(&mystruct, r)
       ...
      })
    r.Methods("GET").Path("/api/{user}/{id}").HandlerFunc(func(w http.ResponseWriter, r *http.Request){

       var mystruct param
       err := b.Bind(&mystruct, r)
       ...

      }
    r.Methods("GET").Path("/api?is_admin=true&test=2").HandlerFunc(func(w http.ResponseWriter, r *http.Request){

       var mystruct query
       err := b.Bind(&mystruct, r)
       ...
      }

      ...
  }
```
