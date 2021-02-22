package http

// import (
// 	"bytes"
// 	"encoding/json"
// 	"net/http"
// 	"net/http/httptest"
// 	"net/url"
// 	"testing"
// )

// func TestPostToken(t *testing.T) {
// 	usecase := NewMockUsecase()
// 	delivery := NewHTTPDelivery(usecase)

// 	data := url.Values{
// 		"username": {"test"},
// 		"password": {"test"},
// 	}
// 	req, err := http.NewRequest("POST", "/v1/token", bytes.NewBufferString(data.Encode()))
// 	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	rr := httptest.NewRecorder()
// 	r := http.NewServeMux()
// 	r.HandleFunc("/v1/token", delivery.routeToken)
// 	r.ServeHTTP(rr, req)

// 	if status := rr.Code; status != http.StatusOK {
// 		t.Errorf("handler returned wrong status code: got %v want %v",
// 			status, http.StatusOK)
// 	}

// 	responsedMap := map[string]interface{}{}
// 	json.Unmarshal(rr.Body.Bytes(), &responsedMap)
// 	t.Logf("got response %v", rr.Body.String())

// 	expected := "bearer"
// 	if responsedMap["token_type"] != expected {
// 		t.Errorf("handler returned unexpected body, error is not match: got %v want userId %v",
// 			rr.Body.String(), expected)
// 	}
// }
