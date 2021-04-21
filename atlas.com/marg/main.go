package main

import (
	"atlas-marg/configurations"
	"atlas-marg/kafka/consumers"
	"atlas-marg/logger"
	"atlas-marg/rest"
	"atlas-marg/tasks"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	l := logger.CreateLogger()

	config, err := configurations.NewConfigurator(l).GetConfiguration()
	if err != nil {
		l.WithError(err).Fatalf("Retrieving the service configuration.")
	}

	consumers.CreateEventConsumers(l)

	go tasks.Register(tasks.NewRespawn(l, config.RespawnInterval))

	rest.CreateRestService(l)

	// trap sigterm or interrupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)

	// Block until a signal is received.
	sig := <-c
	l.Infoln("Shutting down via signal:", sig)
}
