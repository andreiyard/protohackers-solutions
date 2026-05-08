package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net"
)

type Rqst struct {
	Method string   `json:"method"`
	Number *float64 `json:"number"`
}

type Rsp struct {
	Method string `json:"method"`
	Prime  bool   `json:"prime"`
}

func main() {
	ls, err := net.Listen("tcp", ":9999")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := ls.Accept()
		if err != nil {
			log.Println("err when accepting conn", err)
			continue
		}
		go HandleConn(conn)
	}
}

func HandleConn(c net.Conn) {
	// read until end of the message "\n"
	// parse message from json
	// calculate if prime number
	// marshall json answer
	// write to conn
	r := bufio.NewReader(c)
	for {
		data, err := r.ReadBytes('\n')
		if err != nil {
			log.Println("err when reading request", err, string(data))
			break
		}

		rqst := Rqst{}
		log.Printf("rqst: %s", data)
		err = json.Unmarshal(data, &rqst)
		if err != nil || rqst.Number == nil || rqst.Method != "isPrime" {
			log.Printf("request is not valid \"%s\" \n", data[:len(data)-1])
			c.Write([]byte("{}\n"))
			c.Close()
			break
		}

		n := int(*rqst.Number)
		prime := isPrime(n)

		rsp := Rsp{rqst.Method, prime}
		msg, err := json.Marshal(rsp)
		log.Printf("rsp: %s\n", msg)
		if err != nil {
			panic(err)
		}
		c.Write(msg)
		c.Write([]byte("\n"))
	}

}

func isPrime(n int) bool {
	if n < 2 {
		return false
	}
	for i := 2; i*i <= n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}
