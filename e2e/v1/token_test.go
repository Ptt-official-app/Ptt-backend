package e2e

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

func TestPostToken(t *testing.T) {
	data := url.Values{
		"username": {"test"},
		"password": {"test"},
	}
	req, err := http.NewRequest("POST", "https://pttapp.cc/v1/token", bytes.NewBufferString(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()  

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	responsedMap := map[string]interface{}{}
	json.Unmarshal(body, &responsedMap)
	t.Logf("got response %v", body)

	if _, ok := responsedMap["error"]; !ok {
		t.Fatal("no error")
	}
}
