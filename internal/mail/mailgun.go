package mail

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"
)

const (
	scheme     = "https"
	apiVersion = "v3"
	endpoint   = "messages"
)

var (
	_ Mail = &mailgunProvider{}
)

type mailgunProvider struct {
	domain  string
	apiKey  string
	hostURL *url.URL
	client  *http.Client
}

func newMailgunProvider(u *url.URL) *mailgunProvider {
	query := u.Query()
	domain := query.Get("domain")
	apiKey := query.Get("api_key")

	hostURL := &url.URL{
		Scheme: scheme,
		Host:   u.Host,
		Path:   path.Join(apiVersion, domain, endpoint),
	}

	return &mailgunProvider{
		domain:  domain,
		apiKey:  apiKey,
		hostURL: hostURL,
		client:  http.DefaultClient,
	}
}

func (m *mailgunProvider) Send(from, to, title string, body []byte) error {
	formData := url.Values{}

	formData.Add("from", from)
	formData.Add("to", to)
	formData.Add("subject", title)
	formData.Add("text", string(body))

	req, err := http.NewRequest(http.MethodPost, m.hostURL.String(), strings.NewReader(formData.Encode()))
	if err != nil {
		return err
	}

	req.SetBasicAuth("api", m.apiKey)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := m.client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		bodyData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("response status: %s, body: %s", resp.Status, string(bodyData))
	}

	return nil
}
