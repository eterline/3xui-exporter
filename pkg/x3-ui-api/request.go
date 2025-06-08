package x3uiapi

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

type xuiRequest struct {
	xui *XUIClient
	req *http.Request
	res *http.Response
	Ok  bool
}

func (xireq *xuiRequest) get() (code int, err error) {

	xireq.req.Method = "GET"
	response, err := xireq.xui.httpClient.Do(xireq.req)
	if err != nil {
		xireq.Ok = false
		return 0, fmt.Errorf("get request error: %w", err)
	}

	xireq.res = response
	xireq.Ok = true

	return response.StatusCode, nil
}

func (xireq *xuiRequest) post(body io.ReadCloser, json bool) (code int, err error) {

	xireq.req.Method = "POST"
	xireq.req.Body = body

	if json {
		xireq.req.Header.Add("Content-Type", "application/json")
	}

	response, err := xireq.xui.httpClient.Do(xireq.req)
	if err != nil {
		xireq.Ok = false
		return 0, fmt.Errorf("get request error: %w", err)
	}

	xireq.res = response
	xireq.Ok = true

	return response.StatusCode, nil
}

func (xireq *xuiRequest) resolve(v any) error {

	if !xireq.Ok || xireq.res == nil {
		return errors.New("nil response pointer")
	}

	defer xireq.res.Body.Close()
	return json.NewDecoder(xireq.res.Body).Decode(v)
}

func xuiUrl(api, sub string) (*url.URL, error) {
	u, err := url.Parse(api)
	if err != nil {
		return nil, err
	}

	if sub != "" {
		u = u.JoinPath(sub)
	}

	return u, nil
}

func newLoginForm(user, password string) url.Values {
	return map[string][]string{
		"username": {user},
		"password": {password},
	}
}

func caPool(cfg *tls.Config, ca ...string) error {

	caPool := x509.NewCertPool()

	for _, file := range ca {
		data, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read CA certificate: %w", err)
		}

		if !caPool.AppendCertsFromPEM(data) {
			return fmt.Errorf("failed to add CA to pool")
		}
	}

	cfg.RootCAs = caPool
	return nil
}
