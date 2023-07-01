package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"grpcTests/internal/client"
	"grpcTests/internal/network"
	"grpcTests/internal/server"
	"sync"
)

type App struct {
	srv    *server.Server
	client *client.Client
}

func (a *App) Run() {
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
			log.Errorln("Error sending message")
		}
	}
}

func NewApp() (*App, error) {
	srv := server.NewServer()

	// Create own "room" think that is private user
	port := network.RandomPort()
	go func(p string) {
		if err := srv.Run(p); err != nil {
			log.Errorln("Error serving", err)
		}
	}(port)

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

	app.Run()
}
