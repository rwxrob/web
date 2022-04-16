package web_test

import (
	"fmt"
	"net/http"
	ht "net/http/httptest"

	web "github.com/rwxrob/web/pkg"
)

func ExampleFetch_get() {

	// serve get
	handler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, `{"get":"WORKED"}`)
		})
	svr := ht.NewServer(handler)
	defer svr.Close()

	data := map[string]any{}

	req := &web.Req{U: svr.URL, D: data}
	if err := req.Submit(); err != nil {
		fmt.Println(err)
	}
	fmt.Println(req.R.Request.Method)
	fmt.Println(req.R.Status)
	fmt.Println(data["get"])

	// Output:
	// GET
	// 200 OK
	// WORKED
}

func ExampleFetch_post() {

	// serve post
	handler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, `{"post":"WORKED"}`)
		})
	svr := ht.NewServer(handler)
	defer svr.Close()

	data := map[string]any{}

	req := &web.Req{M: `POST`, U: svr.URL, D: data}
	if err := req.Submit(); err != nil {
		fmt.Println(err)
	}
	fmt.Println(req.R.Request.Method)
	fmt.Println(req.R.Status)
	fmt.Println(data["post"])

	// Output:
	// POST
	// 200 OK
	// WORKED
}
