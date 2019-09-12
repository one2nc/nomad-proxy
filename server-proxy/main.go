package main

import (
	"io"
	"log"
	"net"
)

func handler(conn net.Conn) {
	buf := make([]byte, 1)
	if _, err := conn.Read(buf); err != nil {
		if err != io.EOF {
			log.Println("failed to read first RPC byte", "error", err)
		}
		return
	}

	// Means its RPC with tls
	if buf[0] != 0x04 {
		log.Println("NON TLS connection found")
		return
	}

	conn = WrapTls(conn)

	if conn == nil {
		return
	}

	DefaultHandle(conn)
}

func main() {
	RunServer(handler)
}
