package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"time"
)

type TlsParams struct {
	CertFile   string
	KeyFile    string
	CACertFile string
	HardFail   bool
	ServerName string
	SkipTls    bool
}

var HardFail = false

// Return a tls.Config but overlaps with ServerTlsConfig.
// This is a newer function and ServerTlsConfig remains for backward
// compatibility.
func TlsConfig(c *TlsParams) (*tls.Config, error) {
	if c == nil {
		return nil, fmt.Errorf("TlsParams is empty")
	}

	HardFail = c.HardFail
	cfg, err := ServerTlsConfig(c.CertFile, c.KeyFile, c.CACertFile)
	if err != nil {
		return cfg, err
	}

	if c.ServerName != "" {
		cfg.ServerName = c.ServerName
	}

	return cfg, err
}

func certValidator(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
	cert, err := x509.ParseCertificate(rawCerts[0])
	if err != nil {
		return err
	}

	//revoked, ok := revoke.VerifyCertificate(cert)
	var revoked, ok bool
	if !time.Now().Before(cert.NotAfter) || !time.Now().After(cert.NotBefore) {
		revoked, ok = true, true
	}

	if !ok && HardFail {
		return fmt.Errorf("cannot find Cert. Proceed at your own Risk")
	} else if revoked {
		return fmt.Errorf("certificate is revoked")
	}

	return nil
}

// Deprecated: Use TlsConfig instead.
func ServerTlsConfig(cert, key, root string) (*tls.Config, error) {
	certificate, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		return nil, err
	}

	// Create a tlsConfig
	tlsConfig := &tls.Config{
		ServerName:            "tessellate-server",
		ClientAuth:            tls.RequireAndVerifyClientCert,
		Certificates:          []tls.Certificate{certificate},
		VerifyPeerCertificate: certValidator,
	}

	// If rootCert is provided, set ClientCA
	if root != "" {
		ca, err := MakeCertPool(root)
		if err != nil {
			return nil, err
		}
		tlsConfig.ClientCAs = ca
	}

	return tlsConfig, nil
}
