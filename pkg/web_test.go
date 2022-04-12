package web_test

import (
	"fmt"
	"net/http"
	ht "net/http/httptest"

	"github.com/rwxrob/json"
	web "github.com/rwxrob/web/pkg"
)

func ExampleFetch() {

	// serve get
	handler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, `{"get":"t"}`)
		})
	svr := ht.NewServer(handler)
	defer svr.Close()

	// serve get int
	handler0 := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, `20220322075441`)
		})
	svr0 := ht.NewServer(handler0)
	defer svr0.Close()

	// serve post
	handler1 := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, `{"post":"t","c":"t"}`)
		})
	svr1 := ht.NewServer(handler1)
	defer svr1.Close()

	// serve put
	handler2 := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, `{"put":"t"}`)
		})
	svr2 := ht.NewServer(handler2)
	defer svr2.Close()

	// serve patch
	handler3 := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, `{"patch":"t"}`)
		})
	svr3 := ht.NewServer(handler3)
	defer svr3.Close()

	// serve delete
	handler4 := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, `{"delete":"t"}`)
		})
	svr4 := ht.NewServer(handler4)
	defer svr4.Close()

	web.TimeOut = 4

	// create the struct type matching the REST query JSON
	type Data struct {
		Get     string `json:"get"`
		Post    string `json:"post"`
		Put     string `json:"put"`
		Patch   string `json:"patch"`
		Delete  string `json:"delete"`
		Changed string `json:"c"`
		Ignored string `json:"i"`
	}

	data := &Data{
		Changed: "o",
		Ignored: "i",
	}
	jsdata := json.This{data}
	jsdata.Print()

	req := &web.Req{URL: svr.URL, Into: data}

	if err := req.Submit(); err != nil {
		fmt.Println(err)
	}
	jsdata.Print()

	anint := 0
	req = &web.Req{URL: svr0.URL, Into: &anint}
	if err := req.Submit(); err != nil {
		fmt.Println(err)
	}
	fmt.Println(anint)

	req = &web.Req{Method: `POST`, URL: svr1.URL, Into: data}
	if err := req.Submit(); err != nil {
		fmt.Println(err)
	}
	jsdata.Print()

	req = &web.Req{Method: `PUT`, URL: svr2.URL, Into: data}
	if err := req.Submit(); err != nil {
		fmt.Println(err)
	}
	jsdata.Print()

	req = &web.Req{Method: `PATCH`, URL: svr3.URL, Into: data}
	if err := req.Submit(); err != nil {
		fmt.Println(err)
	}
	jsdata.Print()

	req = &web.Req{Method: `DELETE`, URL: svr4.URL, Into: data}
	if err := req.Submit(); err != nil {
		fmt.Println(err)
	}
	jsdata.Print()

	// Output:
	// {"get":"","post":"","put":"","patch":"","delete":"","c":"o","i":"i"}
	// {"get":"t","post":"","put":"","patch":"","delete":"","c":"o","i":"i"}
	// 20220322075441
	// {"get":"t","post":"t","put":"","patch":"","delete":"","c":"t","i":"i"}
	// {"get":"t","post":"t","put":"t","patch":"","delete":"","c":"t","i":"i"}
	// {"get":"t","post":"t","put":"t","patch":"t","delete":"","c":"t","i":"i"}
	// {"get":"t","post":"t","put":"t","patch":"t","delete":"t","c":"t","i":"i"}
}
