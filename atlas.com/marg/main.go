package main

import (
	"atlas-marg/configurations"
	"atlas-marg/kafka/consumers"
	"atlas-marg/logger"
	"atlas-marg/rest"
	"atlas-marg/tasks"
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	l := logger.CreateLogger()
	l.Infoln("Starting main service.")

	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	config, err := configurations.NewConfigurator(l).GetConfiguration()
	if err != nil {
		l.WithError(err).Fatalf("Retrieving the service configuration.")
	}

	consumers.CreateEventConsumers(l, ctx, wg)

	go tasks.Register(tasks.NewRespawn(l, config.RespawnInterval))

	rest.CreateRestService(l, ctx, wg)

	// trap sigterm or interrupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)

	// Block until a signal is received.
	sig := <-c
	l.Infof("Initiating shutdown with signal %s.", sig)
	cancel()
	wg.Wait()
	l.Infoln("Service shutdown.")
}
