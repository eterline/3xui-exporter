package x3uiapi

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync/atomic"
	"time"
)

type X3UIClient struct {
	api        *url.URL
	form       url.Values
	cookie     *atomic.Pointer[http.Cookie]
	httpClient *http.Client
}

func NewClient(api, sub, user, password string, caFile string) (*X3UIClient, error) {

	urlApi, err := xuiUrl(api, sub)
	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{}
	tlsConfig.InsecureSkipVerify = (caFile == "")

	if !tlsConfig.InsecureSkipVerify {
		err := caPool(tlsConfig, caFile)
		if err != nil {
			return nil, err
		}
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
		Timeout: 15 * time.Second,
	}

	formData := newLoginForm(user, password)
	storedCookie := &atomic.Pointer[http.Cookie]{}
	storedCookie.Store(nil)

	cl := &X3UIClient{
		api:        urlApi,
		form:       formData,
		cookie:     storedCookie,
		httpClient: client,
	}

	return cl, nil
}

func (xc *X3UIClient) swapCookie(cookie *http.Cookie) bool {

	if cookie == nil {
		return true
	}

	xc.cookie.Store(cookie)
	return false
}

func (xc *X3UIClient) currentCookie() *http.Cookie {
	return xc.cookie.Load()
}

func (xc *X3UIClient) cookieIsExpired() bool {
	c := xc.currentCookie()
	if c == nil {
		return true
	}

	nowIs := time.Now()
	return c.Expires.After(nowIs)
}

func (xc *X3UIClient) requestLogin(ctx context.Context) error {

	postForm := xc.form.Encode()
	formRd := strings.NewReader(postForm)

	request, err := http.NewRequestWithContext(
		ctx, "POST",
		xc.api.JoinPath("login").String(),
		formRd,
	)

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if err != nil {
		return err
	}

	resp, err := xc.httpClient.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	for _, ce := range resp.Cookies() {
		if ce.Name == "3x-ui" {
			xc.swapCookie(ce)
		}
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("login failed: %s", string(body))
	}

	return nil
}

func (xc *X3UIClient) newRequest(ctx context.Context, path ...string) (*xuiRequest, error) {

	if xc.cookieIsExpired() {
		err := xc.requestLogin(ctx)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, "GET", xc.api.JoinPath(path...).String(), nil)
	if err != nil {
		return nil, err
	}

	req.AddCookie(xc.currentCookie())

	rq := &xuiRequest{
		xui: xc,
		req: req,
		res: nil,
		Ok:  false,
	}

	return rq, nil
}

func (xc *X3UIClient) Inbounds(ctx context.Context) ([]Inbound, error) {

	data := WrapAPI[[]Inbound]{}
	req, err := xc.newRequest(ctx, "panel", "api", "inbounds", "list")
	if err != nil {
		return data.Object, err
	}

	code, err := req.get()
	if err != nil {
		return data.Object, err
	}

	if code > 299 || code < 199 {
		return data.Object, fmt.Errorf("bad status code: %d", code)
	}

	if err := req.resolve(&data); err != nil {
		return data.Object, err
	}

	if !data.Success {
		return data.Object, errors.New(data.Message)
	}

	return data.Object, nil
}

func (xc *X3UIClient) Online(ctx context.Context) (Online, error) {

	data := WrapAPI[Online]{}
	req, err := xc.newRequest(ctx, "panel", "api", "inbounds", "onlines")
	if err != nil {
		return data.Object, err
	}

	code, err := req.post(nil, false)
	if err != nil {
		return data.Object, err
	}

	if code > 299 || code < 199 {
		return data.Object, fmt.Errorf("bad status code: %d", code)
	}

	if err := req.resolve(&data); err != nil {
		return data.Object, err
	}

	if !data.Success {
		return data.Object, errors.New(data.Message)
	}

	return data.Object, nil
}
