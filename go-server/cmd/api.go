package cmd

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"

	"gorm.io/gorm"
)

type Server struct {
	DB      *gorm.DB
	Clients map[string]*Client
	m       sync.RWMutex
	wg      sync.WaitGroup
}

func (s *Server) Run() {
	message := make(chan string)
	listener, err := net.Listen("tcp", ":3000")
	if err != nil {
		err := fmt.Errorf("Error from tcp: %v", err)
		log.Fatal(err)
	}

	defer listener.Close()

	for {
		connection, err := listener.Accept()
		if err != nil {
			err := fmt.Errorf("Error from tcp: %v", err)
			log.Fatal(err)
			continue
		}

		go handleConnection(s, connection, message)
		go getAllMessages(s, message, connection)
	}
}

func handleConnection(s *Server, conn net.Conn, message chan string) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	responses := []string{}
	questions := []string{
		"Welcome, Please enter your name: ",
		"Enter your age: ",
		"What do you like?: ",
	}

	for _, question := range questions {
		writer.WriteString(question + "\n")
		writer.Flush()
		response, _ := reader.ReadString('\n')

		responses = append(responses, response)
	}

	welcomeMessage := fmt.Sprintf("Bonjour Mr: %s", responses[0])
	writer.WriteString(welcomeMessage)
	writer.Flush()
	log.Printf("%s has joined the server", strings.TrimSpace(responses[0]))
	nAge, _ := strconv.Atoi(responses[1])

	client := NewClient(conn.RemoteAddr().String(), responses[0], int8(nAge), responses[2])
	msg, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Disconnected:", err)
		delete(s.Clients, client.Id)
		return
	}

	s.m.Lock()
	s.Clients[client.Id] = client
	// Add person to database
	s.m.Unlock()

	writer.WriteString("Message for you: " + msg)
	writer.Flush()
	message <- msg
}
