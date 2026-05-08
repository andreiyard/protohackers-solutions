package main

import (
	"encoding/binary"
	"errors"
	"io"
	"log"
	"net"
)

type Request struct {
	Kind byte
	Num1 int32
	Num2 int32
}

func main() {
	ls, err := net.Listen("tcp", ":9999")
	if err != nil {
		panic(err)
	}
	for {
		c, err := ls.Accept()
		if err != nil {
			log.Println("err conn accept", err)
		}
		go handleConn(c)
	}
}

func handleConn(c net.Conn) {
	defer c.Close()

	db := make(map[int32]int32)
	var request Request

	for {
		// try to read request from the connection and handle it
		err := binary.Read(c, binary.BigEndian, &request)
		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Println("client closed conn")
			} else {
				log.Println("unable to read request", err)
			}
			break
		}
		log.Println("got request", request)

		err = handleRequest(c, request, db)
		if err != nil {
			log.Println("error while request handling", err)
			break
		}
	}
}

func handleRequest(c net.Conn, request Request, db map[int32]int32) error {
	switch request.Kind {
	case 'I': // insert into db
		db[request.Num1] = request.Num2
	case 'Q': // read from db and send response
		//log.Println("db contents:", db)
		res := calculate(db, request.Num1, request.Num2)
		buf := make([]byte, 4)
		binary.BigEndian.PutUint32(buf, uint32(res))
		_, err := c.Write(buf)
		return err
	default:
		return errors.New("client sent unsupported request type")
	}
	return nil
}

func calculate(db map[int32]int32, start, end int32) int32 {
	var sum int
	var n int
	for time, value := range db {
		if time < start || time > end {
			continue
		}
		sum += int(value)
		n += 1
	}
	if n == 0 {
		return 0
	}
	return int32(sum / n)
}
