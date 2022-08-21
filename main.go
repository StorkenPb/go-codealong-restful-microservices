package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/StorkenPb/restful-microservices-codealong/handlers"
)

func main() {

	l := log.New(os.Stdout, "product-api", log.LstdFlags)

	// Create handlers
	products := handlers.NewProducts(l)

	// Attach handler to the server mux
	mux := http.NewServeMux()
	mux.Handle("/", products)

	httpServer := &http.Server{
		Addr: ":9090",
		Handler: mux,
		IdleTimeout: 120 * time.Second,
		ReadTimeout: 1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}


	go func ()  {
		err := httpServer.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	// Set up channel on which to send signal notifications.
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)
	signal.Notify(signalChannel, os.Kill)

	// Block until a signal is recived
	sig := <- signalChannel
	l.Println("Shutdown in 30 sec", sig)

	timeoutContext, _ := context.WithTimeout(context.Background(), 30*time.Second)
	httpServer.Shutdown(timeoutContext)
}