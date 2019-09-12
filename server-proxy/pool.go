package main

import (
	"crypto/x509"
	"io/ioutil"
)

// MakeCertPool generates a CertPool based on rootPath.
func MakeCertPool(rootPath string) (*x509.CertPool, error) {
	certPool := x509.NewCertPool()
	bs, err := ioutil.ReadFile(rootPath)
	if err != nil {
		return nil, err
	}

	ok := certPool.AppendCertsFromPEM(bs)
	if !ok {
		return nil, err
	}

	return certPool, nil
}
