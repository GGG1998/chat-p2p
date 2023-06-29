package server

import (
	"fmt"
	"google.golang.org/grpc"
	"grpcTests/proto"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"
)

type WrapperServer interface {
	Run() error
}

type Server struct {
	proto.UnimplementedChatServer
	clients map[proto.Chat_SendMessageServer]bool
	mu      sync.Mutex
	srv     *grpc.Server
}

func (s *Server) broadcastMessage(msg *proto.ChatMessage) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for client := range s.clients {
		if err := client.Send(msg); err != nil {
			log.Printf("Błąd podczas wysyłania wiadomości do klienta: %v", err)
		}
	}
}

func (s *Server) SendMessage(stream proto.Chat_SendMessageServer) error {
	s.mu.Lock()
	s.clients[stream] = true
	s.mu.Unlock()

	for {
		msg, err := stream.Recv()
		if err != nil {
			log.Printf("Error during recv from client: %v", err)
			break
		}

		s.broadcastMessage(msg)
	}

	s.mu.Lock()
	delete(s.clients, stream)
	s.mu.Unlock()

	return nil
}

func (s *Server) Run() error {
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)
	port := fmt.Sprintf(":800%d", random.Intn(9))
	socket, err := net.Listen("tcp", port)

	fmt.Printf("\nYour port %s \n", port)

	if err != nil {
		panic(err)
	}

	rpc := grpc.NewServer()
	proto.RegisterChatServer(rpc, s)
	if socket == nil {
		return fmt.Errorf("socket is nil")
	}
	err = rpc.Serve(socket)
	if err != nil {
		return err
	}
	s.srv = rpc
	return nil
}

func NewServer() *Server {
	return &Server{
		clients: make(map[proto.Chat_SendMessageServer]bool),
	}
}
