package pve

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type proxmoxRequests struct {
	API        *url.URL
	TokenID    string
	Token      string
	httpClient *http.Client
}

// newProxmoxRequests - creates proxmox requset api instance
func newProxmoxRequests(api string, tokenID, token, caFile string) (*proxmoxRequests, error) {

	tlsConfig := &tls.Config{}
	tlsConfig.InsecureSkipVerify = (caFile == "")

	if !tlsConfig.InsecureSkipVerify {
		pool, err := caPool(caFile)
		if err != nil {
			return nil, err
		}

		tlsConfig.RootCAs = pool
	}

	uri, err := pveApi(api)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
		Timeout: 15 * time.Second,
	}

	uri = uri.JoinPath("api2", "json")

	reqs := &proxmoxRequests{
		API:        uri,
		TokenID:    tokenID,
		Token:      token,
		httpClient: client,
	}

	return reqs, nil
}

func (pr *proxmoxRequests) joinPath(path ...string) *url.URL {
	return pr.API.JoinPath(path...)
}

func (pr *proxmoxRequests) authTokenString() string {
	return fmt.Sprintf("PVEAPIToken=%s=%s", pr.TokenID, pr.Token)
}

type proxmoxRequest struct {
	httpRequest  *http.Request
	httpResponse *http.Response
	httpClient   *http.Client
	Ok           bool
}

func (pr *proxmoxRequests) request(ctx context.Context, path ...string) *proxmoxRequest {

	request, _ := http.NewRequestWithContext(
		ctx, "", pr.joinPath(path...).String(), nil,
	)

	request.Header.Add("Authorization", pr.authTokenString())

	return &proxmoxRequest{
		httpRequest:  request,
		httpResponse: nil,
		httpClient:   pr.httpClient,
		Ok:           false,
	}
}

func (rp *proxmoxRequest) get() (code int, err error) {

	rp.httpRequest.Method = "GET"
	response, err := rp.httpClient.Do(rp.httpRequest)
	if err != nil {
		rp.Ok = false
		return 0, fmt.Errorf("get request error: %w", err)
	}

	rp.httpResponse = response
	rp.Ok = true

	return response.StatusCode, nil
}

func (rp *proxmoxRequest) resolve(v any) error {

	if !rp.Ok || rp.httpResponse == nil {
		return errors.New("nil response pointer")
	}

	defer rp.httpResponse.Body.Close()
	return json.NewDecoder(rp.httpResponse.Body).Decode(v)
}

func (rp *proxmoxRequest) jsonString() string {

	json := json.RawMessage{}

	if err := rp.resolve(&json); err != nil {
		return ""
	}

	data, err := json.MarshalJSON()
	if err != nil {
		return ""
	}

	return string(data)
}
