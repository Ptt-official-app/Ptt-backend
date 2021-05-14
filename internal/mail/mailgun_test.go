package mail

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path"
	"testing"
)

const (
	testApiKey = "key"
	testDomain = "pttapp.cc"
)

var (
	testQuery = url.Values{"api_key": {testApiKey}, "domain": {testDomain}}
)

func TestMailgunSend_ReturnNoError(t *testing.T) {
	expectPath := path.Join("/", apiVersion, testDomain, endpoint)
	testFrom := "test@example.com"
	testTo := "test"
	testTitle := "test"
	testText := []byte("test")

	testServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != expectPath {
			http.Error(w, "URL path is invalid.", http.StatusBadRequest)
		}

		if r.Method != http.MethodPost {
			http.Error(w, "Request method must be POST", http.StatusBadRequest)
			return
		}

		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}

		if reqFrom := r.FormValue("from"); reqFrom != testFrom {
			http.Error(w, fmt.Sprintf(`FormValue("from"): %q does not match %q`, reqFrom, testFrom), http.StatusBadRequest)
			return
		}
		if reqTo := r.FormValue("to"); reqTo != testTo {
			http.Error(w, fmt.Sprintf(`FormValue("to"): %q does not match %q`, reqTo, testTo), http.StatusBadRequest)
			return
		}
		if reqTitle := r.FormValue("subject"); reqTitle != testTitle {
			http.Error(w, fmt.Sprintf(`FormValue("subject"): %q does not match %q`, reqTitle, testTitle), http.StatusBadRequest)
			return
		}
		if reqText := r.FormValue("text"); reqText != string(testText) {
			http.Error(w, fmt.Sprintf(`FormValue("text"): %q does not match %q`, reqText, string(testText)), http.StatusBadRequest)
			return
		}

	}))
	defer testServer.Close()

	u, err := url.Parse(testServer.URL)
	if err != nil {
		t.Errorf("parse URL error: %v", err)
		return
	}
	u.RawQuery = testQuery.Encode()

	mp := newMailgunProvider(u)
	mp.client = testServer.Client()
	err = mp.Send(testFrom, testTo, testTitle, testText)
	if err != nil {
		t.Errorf("expect nil, got %v", err)
	}
}

func TestMailgunSend_ReturnError(t *testing.T) {
	testServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusBadGateway), http.StatusBadGateway)
		return
	}))
	defer testServer.Close()

	u, err := url.Parse(testServer.URL)
	if err != nil {
		t.Errorf("parse URL error: %v", err)
		return
	}
	u.RawQuery = testQuery.Encode()

	mp := newMailgunProvider(u)
	mp.client = testServer.Client()
	err = mp.Send("test", "test", "test", []byte("test"))
	if err == nil {
		t.Errorf("expect error, got %v", err)
	}
}
