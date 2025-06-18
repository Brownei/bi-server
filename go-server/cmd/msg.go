package cmd

import (
	"fmt"
	"log"
	"net"
	"strings"
)

type Message struct {
	cmd  string
	conn net.Conn
}

func getAllMessages(s *Server, message <-chan string, conn net.Conn) {
	for msg := range message {
		trimmed_msg := strings.TrimSpace(msg)
		fmt.Printf("Messages here: %s\n", trimmed_msg)
		command := strings.Split(trimmed_msg, " ")

		if trimmed_msg == "/exit" {
			curentCleint := s.Clients[conn.RemoteAddr().String()]
			leavingMsg := fmt.Sprintf("Adios Mr %s", strings.TrimSpace(curentCleint.Name))
			log.Printf("%s has left the server", strings.TrimSpace(curentCleint.Name))
			conn.Write([]byte(leavingMsg))
			conn.Close()
			break
		} else if command[0] == "/msg" {
			if len(command) < 3 {
				conn.Write([]byte("You must put the persons name and the message you wanna send"))
			} else {
				curentCleint := s.Clients[conn.RemoteAddr().String()]
				log.Printf("%s wants to message someone", strings.TrimSpace(curentCleint.Name))
			}
		}

		fmt.Print("Done\n")
	}
}
