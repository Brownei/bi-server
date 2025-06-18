package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"

	"github.com/brownei/fast-server/db"
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
	var existingUser db.User
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	client := &Client{}

	responses := []string{}
	questions := []string{
		"Enter your age: ",
		"What do you like?: ",
	}

	writer.WriteString("Welcome, Please enter your name: " + "\n")
	writer.Flush()
	nameResponse, _ := reader.ReadString('\n')

	result := s.DB.First(&existingUser, "name = ?", nameResponse)

	fmt.Print(errors.Is(result.Error, gorm.ErrRecordNotFound))

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		for _, question := range questions {
			writer.WriteString(question + "\n")
			writer.Flush()
			response, _ := reader.ReadString('\n')

			responses = append(responses, response)
		}

		welcomeMessage := fmt.Sprintf("Bonjour Mr: %s", nameResponse)
		writer.WriteString(welcomeMessage)
		writer.Flush()
		log.Printf("%s has joined the server", strings.TrimSpace(nameResponse))
		nAge, _ := strconv.Atoi(responses[0])

		client = NewClient(conn.RemoteAddr().String(), nameResponse, int8(nAge), responses[1])

		s.DB.Create(&db.User{
			Name:  nameResponse,
			Likes: responses[1],
			Age:   uint16(nAge),
		})

	} else {
		welcomeMessage := fmt.Sprintf("Bonjour Mr: %s", nameResponse)
		writer.WriteString(welcomeMessage)
		writer.Flush()
		log.Printf("%s has joined the server", strings.TrimSpace(nameResponse))

		client = NewClient(conn.RemoteAddr().String(), nameResponse, int8(existingUser.Age), existingUser.Likes)

	}

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
