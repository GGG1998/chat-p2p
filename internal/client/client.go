package client

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"grpcTests/proto"
	"log"
)

type WrapperClient interface {
	Connect(string) error
}

type Client struct {
	username string
	client   proto.ChatClient
	stream   proto.Chat_SendMessageClient
	Prompt   *Prompt
}

func (c *Client) SendMessage(message string) error {
	err := c.stream.Send(&proto.ChatMessage{
		Username:    c.username,
		MessageBody: message,
	})
	if err != nil {
		return fmt.Errorf("sending error: %v", err)
	}

	return nil
}

func (c *Client) ReceiveMessages() {
	for {
		msg, err := c.stream.Recv()
		if err != nil {
			break
		}

		log.Printf("[%s] %s\n", msg.Username, msg.MessageBody)
	}
}

func (c *Client) Connect(serverID string) error {
	conn, err := grpc.Dial(serverID, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()

	c.client = proto.NewChatClient(conn)

	c.stream, err = c.client.SendMessage(context.Background())
	if err != nil {
		return fmt.Errorf("create stream error: %v", err)
	}

	return nil
}

func NewClient() *Client {
	prompt := NewPrompt()

	username, err := prompt.Write("username")
	if err != nil {
		username = "anon"
	}

	return &Client{
		Prompt:   prompt,
		username: username,
	}
}
