package main

import (
	"context"
	"fmt"
	"grpcTests/internal/client"
	"grpcTests/internal/server"
	"log"
	"sync"
)

type App struct {
	srv    *server.Server
	client *client.Client
}

func (a *App) Run(ctx context.Context) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		a.client.ReceiveMessages()
	}()

	// main program
	for {
		msg, _ := a.client.Prompt.Write("")
		if err := a.client.SendMessage(msg); err != nil {
			log.Printf("Error sending message")
		}
	}
}

func NewApp() (*App, error) {
	srv := server.NewServer()

	// Create own "room" think that is private user
	go func() {
		if err := srv.Run(); err != nil {
			fmt.Println("Error serving", err)
		}
	}()

	serverId, err := client.NewPrompt().Write("serverID")
	if err != nil {
		return nil, fmt.Errorf("prompt failed: %v", err)
	}

	clt := client.NewClient()
	if err := clt.Connect(serverId); err != nil {
		return nil, fmt.Errorf("client error")
	}

	return &App{
		srv:    srv,
		client: clt,
	}, nil
}

func main() {
	app, err := NewApp()
	if err != nil {
		_ = fmt.Errorf("app error: %w", err)
	}

	app.Run(context.Background())
}
