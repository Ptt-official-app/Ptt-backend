package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestPostToken(t *testing.T) {
	usecase := &MockUsecase{}
	delivery := NewHTTPDelivery(usecase)

	data := url.Values{
		"username": {"test"},
		"password": {"test"},
	}
	req, err := http.NewRequest("POST", "/v1/token", bytes.NewBufferString(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.RemoteAddr = "203.0.113.9:4567"
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	r := http.NewServeMux()
	r.HandleFunc("/v1/token", delivery.routeToken)
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	responsedMap := map[string]interface{}{}
	err = json.Unmarshal(rr.Body.Bytes(), &responsedMap)
	if err != nil {
		t.Errorf("get unexpect json: %v", err)
	}
	t.Logf("got response %v", rr.Body.String())

	expected := "bearer"
	if responsedMap["token_type"] != expected {
		t.Errorf("handler returned unexpected body, error is not match: got %v want userId %v",
			rr.Body.String(), expected)
	}
	if usecase.loginUserID != "test" {
		t.Errorf("recorded login user = %q, want test", usecase.loginUserID)
	}
	if usecase.loginIP != "203.0.113.9" {
		t.Errorf("recorded login IP = %q, want 203.0.113.9", usecase.loginIP)
	}
}

func TestRequestRemoteIP(t *testing.T) {
	testCases := []struct {
		remoteAddr string
		want       string
	}{
		{remoteAddr: "127.0.0.1:8080", want: "127.0.0.1"},
		{remoteAddr: "[2001:db8::1]:8080", want: "2001:db8::1"},
		{remoteAddr: "2001:db8::1", want: "2001:db8::1"},
	}

	for _, testCase := range testCases {
		req := httptest.NewRequest(http.MethodPost, "/v1/token", nil)
		req.RemoteAddr = testCase.remoteAddr
		if got := requestRemoteIP(req); got != testCase.want {
			t.Errorf("requestRemoteIP(%q) = %q, want %q", testCase.remoteAddr, got, testCase.want)
		}
	}
}

func TestPostTokenFormEmpty(t *testing.T) {
	usecase := NewMockUsecase()
	delivery := NewHTTPDelivery(usecase)

	data := url.Values{
		"username": {},
		"password": {},
	}
	req, err := http.NewRequest("POST", "/v1/token", bytes.NewBufferString(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	r := http.NewServeMux()
	r.HandleFunc("/v1/token", delivery.routeToken)
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	responsedMap := map[string]interface{}{}
	err = json.Unmarshal(rr.Body.Bytes(), &responsedMap)
	if err != nil {
		t.Errorf("get unexpect json: %v", err)
	}
	t.Logf("got response %v", rr.Body.String())

	expected := "empty username"
	if responsedMap["error_description"] != expected {
		t.Errorf("handler returned unexpected body, error is not match: got %v want userId %v",
			rr.Body.String(), expected)
	}
}
