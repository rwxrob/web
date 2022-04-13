// Copyright 2022 web Robert Muhlestein.
// SPDX-License-Identifier: Apache-2.0

// Package web provides high-level functions that are called from the Go
// Bonzai branch of the same name providing universal access to the core
// functionality.
package web

import (
	"context"
	"encoding"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	rwxjson "github.com/rwxrob/json"
	"gopkg.in/yaml.v3"
)

// TimeOut is a package global timeout for any of the high-level https
// query functions in this package. The default value is 60 seconds.
var TimeOut int = 60

// Context is convenient alias to context.Context interface.
type Context context.Context

// HTTPError is an error for anything other than an HTTP response in the
// 200-299 range including the 300 redirects (which should be handled by
// the Req.Submit successfully before returning). http.Response is
// embedded directly for details.
type HTTPError struct {
	Resp *http.Response
}

// Error fulfills the error interface.
func (e HTTPError) Error() string { return e.Resp.Status }

// ReqSyntaxError is for any error involving the incorrect
// definition of Req fields (such as including a question mark in
// the URL, etc.).
type ReqSyntaxError struct {
	Message string
}

// Error fulfills the error interface.
func (e ReqSyntaxError) Error() string { return e.Message }

// Client provides a way to change the default HTTP client for any
// further package HTTP request function calls. The Client can also be
// set in any Req by assigning the to the field of the same name. By
// default, it is set to http.DefaultClient. This is particularly useful
// when creating mockups and other testing.
var Client = http.DefaultClient

// Headers contains headers to be added to a Req. Unlike the
// specification, only one Header of a give name is allowed. For more
// precision the net/http library directly should be used instead.
type Headers map[string]string

// Req is a human-friendly way to think of web requests. This design
// is closer a pragmatic curl requests than the canonical specification
// (unique headers, for example). The type and parameters of the web
// request and response are derived from the Req fields.
//

//
// The Body can be one of several types that till trigger what is
// submitted as data portion of the request:
//
//     url.Values - will be www-form-encoded and trigger type
//     byte       - uuencoded binary data
//     string     - plain text
//
// Note that Req has no support for multi-part MIME. Use net/http
// directly if such is required.
//
// The Data field can also be any of several types that trigger how the
// received data is handled:
//
//     []byte           - uudecoded binary
//     string           - plain text string
//     io.Writer        - keep as is
//     json.This        - unmarshaled JSON data into This
//     any              - unmarshaled JSON data
//
// Passing the query string as url.Values automatically add
// a question mark (?) followed by the URL encoded values to the end of
// the URL which may present a problem if the URL already has a query
// string. Encouraging the use of url.Values for passing the query
// string serves as a reminder that all query strings should be URL
// encoded (as is often forgotten).
type Req struct {
	Method  string     // GET, POST, (default GET)
	URL     string     // base url, never any query string
	Query   url.Values // query string to append to URL
	Headers Headers    // never more than one of same
	Body    any        // body data, if url.Values will JSON encode
	Data    any        // where to put the response data
	Context Context    // trigger requests with context
	Resp    any        // usually http.Response
}

// Submit synchronously sends the Req to server and populates the
// response from the server into Data. Anything but a response in the
// 200s will result in an HTTPError. See Req for details on how
// inspection of Req will change the behavior of Submit
// automatically. It Req.Context is nil a context.WithTimeout will
// be used and with the value of web.TimeOut.
//
// For convenience, will produce a ReqSyntaxError for any of the
// following conditions:
//
//     * URL contains '?' (use Query instead)
//     * Unsupported Body type
//     * Unsupported Data type
//
func (req *Req) Submit() error {

	if req.Method == "" {
		req.Method = `GET`
	}

	if strings.Index(req.URL, "?") >= 0 {
		return ReqSyntaxError{`URL contains '?' (use Query instead)`}
	}
	req.URL = req.URL + "?" + req.Query.Encode()
	var bodyReader io.Reader

	if req.Headers == nil {
		req.Headers = Headers{}
	}

	var buf string

	switch v := req.Body.(type) {
	case url.Values:
		buf = v.Encode()
		req.Headers["Content-Type"] = "application/x-www-form-urlencoded"
	case []byte:
		log.Println("planned, but unimplemented, would uuencode")
		//req.Headers["Content-Length"] = strconv.Itoa(len(uuencoded))
	case string:
		buf = v
	case yaml.Marshaler:
		byt, err := yaml.Marshal(v)
		if err != nil {
			return err
		}
		buf = string(byt)
	case json.Marshaler:
		byt, err := json.Marshal(v)
		if err != nil {
			return err
		}
		buf = string(byt)
	case encoding.TextMarshaler:
		byt, err := v.MarshalText()
		if err != nil {
			return err
		}
		buf = string(byt)
	case fmt.Stringer:
		buf = v.String()
	default:
		buf = fmt.Sprintf("%v", v)
	}

	bodyReader = strings.NewReader(buf)
	req.Headers["Content-Length"] = strconv.Itoa(len(buf))

	httpreq, err := http.NewRequest(req.Method, req.URL, bodyReader)
	if err != nil {
		return err
	}

	if req.Headers != nil {
		for k, v := range req.Headers {
			httpreq.Header.Add(k, v)
		}
	}

	if req.Context == nil {
		dur := time.Duration(time.Second * time.Duration(TimeOut))
		ctx, cancel := context.WithTimeout(context.Background(), dur)
		defer cancel()
		httpreq = httpreq.WithContext(ctx)
	}

	res, err := Client.Do(httpreq)
	req.Resp = res

	if err != nil {
		return err
	}

	if !(200 <= res.StatusCode && res.StatusCode < 300) {
		return HTTPError{res}
	}

	resbytes, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if len(resbytes) == 0 {
		return nil
	}

	switch req.Data.(type) {
	case map[string]any:
		return yaml.Unmarshal(resbytes, req.Data)
	case string:
		req.Data = string(resbytes)
	case []byte:
		log.Println("planned, but unimplemented, would uuencode")
		// v = uudecode(resbytes)
	case yaml.Unmarshaler:
		return yaml.Unmarshal(resbytes, req.Data)
	case io.Writer:
		return nil
	case rwxjson.This:
		log.Println("rwxjson, planned, but unimplemented")
	default:
		return yaml.Unmarshal(resbytes, req.Data)
	}

	return nil

}
