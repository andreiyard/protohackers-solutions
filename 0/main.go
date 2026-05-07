package main

import (
	"flag"
	"fmt"
	"io"
	"net"
)

func main() {
	ip := flag.String("address", "", "listen IP address")
	port := flag.Int("port", 8080, "listen TCP port")
	flag.Parse()

	address := fmt.Sprint(*ip, ":", *port)

	fmt.Println("Server listening on", address)

	ln, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	fmt.Println("Conn accepted", conn.RemoteAddr().String())
	io.Copy(conn, conn)
}
