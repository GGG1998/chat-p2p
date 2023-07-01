package server

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"grpcTests/proto"
	"net"
	"sync"
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
			log.Printf("Error send: %v", err)
		}
	}
}

func (s *Server) SendMessage(stream proto.Chat_SendMessageServer) error {
	s.mu.Lock()
	s.clients[stream] = true
	s.mu.Unlock()

	var _err error
	for {
		msg, err := stream.Recv()
		if err != nil {
			_err = err
			break
		}

		s.broadcastMessage(msg)
	}

	s.mu.Lock()
	delete(s.clients, stream)
	s.mu.Unlock()

	return _err
}

func (s *Server) Run(port string) error {
	socket, err := net.Listen("tcp", port)

	log.Infof("\nYour port %s \n", port)

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
