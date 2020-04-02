// Copyright 2019 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package opsmngr

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/go-test/deep"
)

var (
	mux *http.ServeMux

	ctx = context.TODO()

	client *Client

	server *httptest.Server
)

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	client = NewClient(nil)
	u, _ := url.Parse(server.URL)
	client.BaseURL = u
}

func teardown() {
	server.Close()
}

func testMethod(t *testing.T, r *http.Request, expected string) {
	if expected != r.Method {
		t.Errorf("Request method = %v, expected %v", r.Method, expected)
	}
}

func testURLParseError(t *testing.T, err error) {
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
	if err, ok := err.(*url.Error); !ok || err.Op != "parse" {
		t.Errorf("Expected URL parse error, got %+v", err)
	}
}

func testClientDefaultBaseURL(t *testing.T, c *Client) {
	if c.BaseURL == nil || c.BaseURL.String() != DefaultBaseURL {
		t.Errorf("NewClient BaseURL = %v, expected %v", c.BaseURL, DefaultBaseURL)
	}
}

func testClientDefaultUserAgent(t *testing.T, c *Client) {
	if c.UserAgent != DefaultUserAgent {
		t.Errorf("NewClient UserAgent = %v, expected %v", c.UserAgent, DefaultUserAgent)
	}
}

func testClientDefaults(t *testing.T, c *Client) {
	testClientDefaultBaseURL(t, c)
	testClientDefaultUserAgent(t, c)
}

func TestNewClient(t *testing.T) {
	c := NewClient(nil)
	testClientDefaults(t, c)
}

func TestNew(t *testing.T) {
	c, err := New(nil)

	if err != nil {
		t.Fatalf("New(): %v", err)
	}
	testClientDefaults(t, c)
}

type testRequestBody struct {
	TestName    string `json:"testName"`
	TestCounter int64  `json:"testCounter"`
	TestData    string `json:"testUserData"`
}

func TestNewRequest_withUserData(t *testing.T) {
	c := NewClient(nil)

	requestPath := "foo"

	inURL, outURL := requestPath, DefaultBaseURL+requestPath
	inBody, outBody := &testRequestBody{TestName: "l", TestData: "u"},
		`{"testName":"l","testCounter":0,`+
			`"testUserData":"u"}`+"\n"
	req, _ := c.NewRequest(ctx, http.MethodGet, inURL, inBody)

	// test relative URL was expanded
	if req.URL.String() != outURL {
		t.Errorf("NewRequest(%v) URL = %v, expected %v", inURL, req.URL, outURL)
	}

	// test body was JSON encoded
	body, _ := ioutil.ReadAll(req.Body)
	if string(body) != outBody {
		t.Errorf("NewRequest(%v)Body = %v, expected %v", inBody, string(body), outBody)
	}

	// test default user-agent is attached to the request
	userAgent := req.Header.Get("User-Agent")
	if c.UserAgent != userAgent {
		t.Errorf("NewRequest() User-Agent = %v, expected %v", userAgent, c.UserAgent)
	}
}

func TestNewRequest_badURL(t *testing.T) {
	c := NewClient(nil)
	_, err := c.NewRequest(ctx, http.MethodGet, ":", nil)
	testURLParseError(t, err)
}

func TestNewRequest_withCustomUserAgent(t *testing.T) {
	ua := "testing/0.0.1"
	c, err := New(nil, SetUserAgent(ua))

	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}

	req, _ := c.NewRequest(ctx, http.MethodGet, "/foo", nil)

	expected := fmt.Sprintf("%s %s", ua, DefaultUserAgent)
	if got := req.Header.Get("User-Agent"); got != expected {
		t.Errorf("New() UserAgent = %s; expected %s", got, expected)
	}
}

func TestClient_Do(t *testing.T) {
	setup()
	defer teardown()

	type foo struct {
		A string
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if m := http.MethodGet; m != r.Method {
			t.Errorf("Request method = %v, expected %v", r.Method, m)
		}
		_, _ = fmt.Fprint(w, `{"A":"a"}`)
	})

	req, _ := client.NewRequest(ctx, http.MethodGet, "/", nil)
	body := new(foo)
	_, err := client.Do(context.Background(), req, body)
	if err != nil {
		t.Fatalf("Do(): %v", err)
	}

	expected := &foo{"a"}
	if !reflect.DeepEqual(body, expected) {
		t.Errorf("Response body = %v, expected %v", body, expected)
	}
}

func TestClient_Do_httpError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Bad Request", 400)
	})

	req, _ := client.NewRequest(ctx, http.MethodGet, "/", nil)
	_, err := client.Do(context.Background(), req, nil)

	if err == nil {
		t.Error("Expected HTTP 400 error.")
	}
}

// Test handling of an error caused by the internal http client's Do()
// function.
func TestClient_Do_redirectLoop(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/", http.StatusFound)
	})

	req, _ := client.NewRequest(ctx, http.MethodGet, "/", nil)
	_, err := client.Do(context.Background(), req, nil)

	if err == nil {
		t.Error("Expected error to be returned.")
	}
	if err, ok := err.(*url.Error); !ok {
		t.Errorf("Expected a URL error; got %#v.", err)
	}
}

func TestClient_OnRequestCompleted(t *testing.T) {
	setup()
	defer teardown()

	type foo struct {
		A string
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if m := http.MethodGet; m != r.Method {
			t.Errorf("Request method = %v, expected %v", r.Method, m)
		}
		_, _ = fmt.Fprint(w, `{"A":"a"}`)
	})

	req, _ := client.NewRequest(ctx, http.MethodGet, "/", nil)
	body := new(foo)
	var completedReq *http.Request
	var completedResp string
	client.OnRequestCompleted(func(req *http.Request, resp *http.Response) {
		completedReq = req
		b, err := httputil.DumpResponse(resp, true)
		if err != nil {
			t.Errorf("Failed to dump response: %s", err)
		}
		completedResp = string(b)
	})
	_, err := client.Do(context.Background(), req, body)
	if err != nil {
		t.Fatalf("Do(): %v", err)
	}
	if diff := deep.Equal(req, completedReq); diff != nil {
		t.Error(diff)
	}
	expected := `{"A":"a"}`
	if !strings.Contains(completedResp, expected) {
		t.Errorf("expected response to contain %v, Response = %v", expected, completedResp)
	}
}

func TestSetUserAgent(t *testing.T) {
	ua := "testing/0.0.1"
	c, err := New(nil, SetUserAgent(ua))

	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}

	expected := fmt.Sprintf("%s %s", ua, DefaultUserAgent)
	if got := c.UserAgent; got != expected {
		t.Errorf("New() UserAgent = %s; expected %s", got, expected)
	}
}

func TestSetBaseURL(t *testing.T) {
	baseURL := "http://localhost/foo"
	c, err := New(nil, SetBaseURL(baseURL))

	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}

	expected := baseURL
	if got := c.BaseURL.String(); got != expected {
		t.Errorf("New() BaseURL = %s; expected %s", got, expected)
	}
}

func TestSetBaseURL_badURL(t *testing.T) {
	baseURL := ":"
	_, err := New(nil, SetBaseURL(baseURL))

	testURLParseError(t, err)
}

func TestOptionSkipVerify(t *testing.T) {
	_, err := New(nil, OptionSkipVerify)

	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}
}

func TestOptionCAValidate(t *testing.T) {
	_, err := New(nil, OptionCAValidate(""))

	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}
}
