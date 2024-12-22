package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"yaly-1/internal/rest"
	"yaly-1/internal/service"
)

func main() {
	log.Println("Hello")

	port := "8080"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}
	server := rest.NewServer(port, service.NewCalcService())

	graceChannel := make(chan os.Signal, 1)
	signal.Notify(graceChannel, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.Run(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-graceChannel

	if err := server.Shutdown(); err != nil {
		log.Fatal(err)
	}
	log.Println("Goodbye")
}
