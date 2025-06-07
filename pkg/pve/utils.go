package pve

import (
	"crypto/x509"
	"fmt"
	"net/url"
	"os"
)

func caPool(ca ...string) (crt *x509.CertPool, err error) {

	caPool := x509.NewCertPool()

	for _, file := range ca {
		data, err := os.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("failed to read CA certificate: %w", err)
		}

		if !caPool.AppendCertsFromPEM(data) {
			return nil, fmt.Errorf("failed to add CA to pool")
		}
	}

	return caPool, nil
}

func pveApi(uri string) (*url.URL, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	u.Scheme = "https"

	return u, nil
}
