package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

func main() {
	db := make(map[string]string)

	address, err := net.ResolveUDPAddr("udp", ":20000")
	conn, err := net.ListenUDP("udp", address)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	buf := make([]byte, 65535) // max datagram size
	for {
		n, remoteAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Println(err)
			continue
		}
		log.Printf("got packet from %s len %d", remoteAddr, n)

		// process request
		request := string(buf[0:n])
		log.Printf("request: %q", request)
		if strings.Contains(request, "=") {
			//set
			parts := strings.SplitN(request, "=", 2)
			key := parts[0]
			value := parts[1]
			if key == "version" {
				continue
			}
			db[key] = value
		} else {
			//get
			key := request
			var value string
			if key == "version" {
				log.Println("current db", db)
				value = "zalopa"
			} else {
				got, ok := db[key]
				value = got
				if !ok {
					value = ""
				}
			}
			response := fmt.Sprintf("%s=%s", key, value)
			log.Printf("responding with %q to %s", response, remoteAddr)
			conn.WriteToUDP([]byte(response), remoteAddr)
		}
	}
}
