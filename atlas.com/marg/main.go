package main

import (
	"atlas-marg/configurations"
	"atlas-marg/kafka/consumers"
	"atlas-marg/logger"
	_map "atlas-marg/map"
	"atlas-marg/rest"
	"atlas-marg/tasks"
	"atlas-marg/tracing"
	"context"
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const serviceName = "atlas-marg"

func main() {
	l := logger.CreateLogger(serviceName)
	l.Infoln("Starting main service.")

	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	tc, err := tracing.InitTracer(l)(serviceName)
	if err != nil {
		l.WithError(err).Fatal("Unable to initialize tracer.")
	}
	defer func(tc io.Closer) {
		err := tc.Close()
		if err != nil {
			l.WithError(err).Errorf("Unable to close tracer.")
		}
	}(tc)

	config, err := configurations.NewConfigurator(l).GetConfiguration()
	if err != nil {
		l.WithError(err).Fatalf("Retrieving the service configuration.")
	}

	consumers.CreateEventConsumers(l, ctx, wg)

	go tasks.Register(tasks.NewRespawn(l, config.RespawnInterval))

	rest.CreateService(l, ctx, wg, "/ms/marg", _map.InitResource)

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
