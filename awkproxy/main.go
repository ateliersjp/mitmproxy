package main

import (
	"os"
	"log"
	"net"

	"github.com/ateliersjp/http"
)

const (
	LISTEN_PROTOCOL      = "unix"
	DIAL_PROTOCOL        = "tcp"
	LISTEN_ADDRESS       = "/var/run/mitmproxy/awkproxy.sock"
	DIAL_ADDRESS         = "localhost:8080"
)

func main() {
	os.Remove(LISTEN_ADDRESS)
	ln, err := net.Listen(LISTEN_PROTOCOL, LISTEN_ADDRESS)
	if err != nil {
		log.Fatal(err)
	}

	err = os.Chmod(LISTEN_ADDRESS, 0666)
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go handle(conn)
	}
}

func handle(conn net.Conn) {
	defer conn.Close()
	if req, err := http.ReadMsg(conn); err == nil {
		if remote, err := net.Dial(DIAL_PROTOCOL, DIAL_ADDRESS); err == nil {
			defer remote.Close()
			if sed := getSed(req); sed != nil {
				wg := NewWaitGroup()
				go wg.Copy(remote, req.Reader())
				if res, err := http.ReadMsg(remote); err == nil {
					res.Headers = sed.TransformHeaders(res.Headers)
					res.Body = sed.TransformBody(res.Body)
					go wg.Copy(conn, res.Reader())
				} else {
					wg.Done()
				}
				wg.Wait()
			}
		}
	}
}
