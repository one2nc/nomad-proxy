package main

import (
	"github.com/tsocial/catoolkit/tlsproxy"
	"net"
	"io"
	"log"
)

func handler(conn net.Conn) {
	buf := make([]byte, 1)
	if _, err := conn.Read(buf); err != nil {
		if err != io.EOF {
			log.Println("failed to read first RPC byte", "error", err)
		}
		conn.Close()
		return
	}

	// Means its RPC with tls
	if buf[0] != 0x04 {
		log.Println("NON TLS connection found")
		conn.Close()
		return
	}

	conn = tlsproxy.WrapTls(conn)

	if conn == nil {
		return
	}

	tlsproxy.DefaultHandle(conn)
}



func main() {
	tlsproxy.RunServer(handler)
}