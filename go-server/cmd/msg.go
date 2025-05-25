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

		if trimmed_msg == "exit" {
			curentCleint := s.Clients[conn.RemoteAddr().String()]
			leavingMsg := fmt.Sprintf("Adios Mr %s", strings.TrimSpace(curentCleint.Name))
			log.Printf("%s has left the server", strings.TrimSpace(curentCleint.Name))
			conn.Write([]byte(leavingMsg))
			conn.Close()
			break
		}

		fmt.Print("Done\n")
	}
}
