package main

import (
	"bufio"
	"fmt"
	"log"
	"maps"
	"net"
	"slices"
	"strings"
)

type ConnWithUsername struct {
	Conn     net.Conn
	Username string
}

type Message struct {
	From net.Conn
	Msg  string
}

type Server struct {
	users     map[net.Conn]string
	Join      chan ConnWithUsername
	Leave     chan net.Conn
	Broadcast chan Message
}

func NewServer() Server {
	return Server{make(map[net.Conn]string), make(chan ConnWithUsername), make(chan net.Conn), make(chan Message)}
}

func (s Server) handleConn(c net.Conn) {
	defer c.Close()

	// get user name and verify it
	fmt.Fprintln(c, "Welcome to budgetchat! What shall I call you?")
	scanner := bufio.NewScanner(c)
	if !scanner.Scan() {
		log.Println("EOF while waiting for username")
		return
	}
	username := strings.TrimSpace(scanner.Text())

	if isValidUsername(username) {
		// send signal to server that user joined
		s.Join <- ConnWithUsername{c, username}
	} else {
		log.Println("user provided bad username", username)
		fmt.Fprintln(c, "Username not valid (should have 0<chars<21 and a-zA-Z0-9)")
		return
	}

	// when user closes connection send leave signal
	defer func() { s.Leave <- c }()

	// process new messages from user in a loop, send to broadcast
	for scanner.Scan() {
		line := scanner.Text()
		s.Broadcast <- Message{c, line}
	}
}

func isValidUsername(name string) bool {
	l := len(name)
	if l < 1 || l > 20 {
		return false
	}
	for _, r := range name {
		if !('a' <= r && r <= 'z') && !('A' <= r && r <= 'Z') && !('0' <= r && r <= '9') {
			return false
		}
	}
	return true
}

func (s *Server) run() {
	for {
		select {
		case conWithName := <-s.Join:
			log.Println("New user joined", conWithName, conWithName.Conn.RemoteAddr())

			// notify all users
			s.doBroadcast(conWithName.Conn, fmt.Sprintf("* %s has entered the room", conWithName.Username), "")

			// show info to new user
			userList := strings.Join(slices.Collect(maps.Values(s.users)), ", ")
			fmt.Fprintf(conWithName.Conn, "* The room contains: %s\n", userList)

			// save
			s.users[conWithName.Conn] = conWithName.Username

		case con := <-s.Leave:
			log.Println("User left", con.RemoteAddr())

			// notfiy all users
			s.doBroadcast(con, fmt.Sprintf("* %s has left the room", s.users[con]), "")

			// remove from memory
			delete(s.users, con)

		case msg := <-s.Broadcast:
			log.Println("New message received", msg, msg.From.RemoteAddr())
			s.doBroadcast(msg.From, msg.Msg, s.users[msg.From])
		}
	}
}

func (s Server) doBroadcast(from net.Conn, msg string, username string) {
	prefix := ""
	if username != "" {
		prefix = fmt.Sprintf("[%s] ", username)
	}

	for c := range s.users {
		if c != from {
			fmt.Fprintf(c, "%s%s\n", prefix, msg)
		}
	}

}

func main() {
	ls, err := net.Listen("tcp", ":9999")
	if err != nil {
		log.Fatalln(err)
	}

	s := NewServer()
	go s.run()

	for {
		c, err := ls.Accept()
		if err != nil {
			log.Println("failed to accept client")
		}
		go s.handleConn(c)
	}
}
