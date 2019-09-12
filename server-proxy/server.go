package main

import (
	"crypto/tls"
	"log"
	"net"

	"fmt"

	"net/http"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	Listen   = kingpin.Flag("listen", "Listening on port").Default("9899").String()
	SkipTLS  = kingpin.Flag("skip-tls", "Skip using TLS while connecting to the server.").Default("false").Bool()
	RootCert = kingpin.Flag("root-cert-file", "Root Cert File").Envar("ROOT_CERT_FILE").String()
	CertFile = kingpin.Flag("cert-file", "Cert File").Envar("CERT_FILE").String()
	KeyFile  = kingpin.Flag("key-file", "Key File").Envar("KEY_FILE").String()

	Upstream = kingpin.Flag("upstream-addr", "Upstream server address, the main server.").
			Envar("UPSTREAM_ADDR").Default("localhost:8080").String()
	TlsServerName = kingpin.Flag("tls-server-name", "ServerName in tls Config").
			Default("localhost").Envar("TLS_SERVER_NAME").String()
)

// Default TlsParams
func cliParams() *TlsParams {
	return &TlsParams{
		CertFile:   *CertFile,
		KeyFile:    *KeyFile,
		CACertFile: *RootCert,
		HardFail:   false,
		SkipTls:    *SkipTLS,
		ServerName: *TlsServerName,
	}
}

func ServeHTTP(listen string, handler http.Handler, p *TlsParams) error {
	if p == nil {
		p = cliParams()
	}

	tlsConfig, err := TlsConfig(p)
	if err != nil {
		return err
	}

	s := http.Server{
		Addr:      listen,
		Handler:   handler,
		TLSConfig: tlsConfig,
	}

	return s.ListenAndServeTLS("", "")
}

func Run(listen string, handle func(conn net.Conn)) {
	RunWithParams(listen, handle, cliParams())
}

func RunWithParams(listen string, handle func(conn net.Conn), p *TlsParams) {
	run(listen, func(conn net.Conn) {
		if !(p.SkipTls) {
			conn = SwitchConnTls(conn, p)
		}
		handle(conn)
	})
}

func run(listen string, handle func(conn net.Conn)) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Printf("Listening on :%v\n", listen)
	ln, err := net.Listen("tcp", ":"+listen)
	if err != nil {
		log.Fatal(err)
		return
	}

	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go func() {
			defer conn.Close()
			handle(conn)
		}()
	}
}

// Deprecated: Use run instead.
// Runs the server at port with the handler passed as a parameter.
func RunServer(handle func(conn net.Conn)) {
	kingpin.Parse()
	log.SetFlags(log.LstdFlags | log.Llongfile)
	run(*Listen, handle)
}

// Deprecated: Use RunServer and pass DefaultHandle instead.
func DefaultTlsHandle(conn net.Conn) {
	DefaultHandle(WrapTls(conn))
}

// Default Handle which manages a TCP Connection.
func DefaultHandle(conn net.Conn) {
	if conn == nil {
		return
	}

	rConn, err := net.Dial("tcp", *Upstream)
	if err != nil {
		log.Printf("error while connecting to upstream: %s. %+v", *Upstream, err)
		return
	}

	defer func() {
		log.Println("Closing remote conn")
		rConn.Close()
	}()

	// Request loop
	go func() {
		for {
			data := make([]byte, 1024*1024)
			n, err := conn.Read(data)
			if err != nil {
				log.Printf("error while reading from client: %+v", err)
				break
			}
			rConn.Write(data[:n])
		}
	}()

	// Response loop
	for {
		data := make([]byte, 1024*1024)
		n, err := rConn.Read(data)
		if err != nil {
			log.Printf("error while reading from remote: %+v", err)
			break
		}
		conn.Write(data[:n])
	}
}

// Switches a basic TCP connection to a TLS connection.
func SwitchConnTls(conn net.Conn, p *TlsParams) net.Conn {
	tlsConfig, err := TlsConfig(p)
	if err != nil {
		panic(err)
	}

	conn = tls.Server(conn, tlsConfig)
	tlsConn, ok := conn.(*tls.Conn)
	if !ok {
		log.Println("expected TLS connection", "got", fmt.Sprintf("%T", conn))
		return nil
	}

	if err := tlsConn.Handshake(); err != nil {
		log.Println("failed TLS handshake", "remote_addr", tlsConn.RemoteAddr(), "error", err)
		return nil
	}

	return tlsConn
}

// Deprecated: Backward compatible wrapper for SwitchConnTls
func WrapTls(conn net.Conn) net.Conn {
	return SwitchConnTls(conn, cliParams())
}
