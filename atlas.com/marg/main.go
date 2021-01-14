package main

import (
	"atlas-marg/configurations"
	"atlas-marg/consumers"
	"atlas-marg/handlers"
	"atlas-marg/tasks"
	"context"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func main() {
	l := log.New(os.Stdout, "marg ", log.LstdFlags | log.Lmicroseconds)
	c, err := configurations.NewConfigurator(l).GetConfiguration()
	if err != nil {
		l.Fatal("[ERROR] Retrieving the service configuration")
	}

	go consumers.NewMapChanged(l, context.Background()).Init()
	go consumers.NewCharacterStatus(l, context.Background()).Init()

	go tasks.Register(tasks.NewRespawn(l, c.RespawnInterval))

	handleRequests(l)
}

func handleRequests(l *log.Logger) {
	//TODO this needs to be updated
	router := mux.NewRouter().StrictSlash(true).PathPrefix("/ms/mrg").Subrouter()
	router.Use(commonHeader)
	router.Handle("/docs", middleware.Redoc(middleware.RedocOpts{BasePath: "/ms/mrg", SpecURL: "/ms/mrg/swagger.yaml"}, nil))
	router.Handle("/swagger.yaml", http.StripPrefix("/ms/mrg", http.FileServer(http.Dir("/"))))

	m := handlers.NewMap(l)
	mRouter := router.PathPrefix("/worlds").Subrouter()
	mRouter.HandleFunc("/{worldId}/channels/{channelId}/maps/{mapId}/characters", m.GetMapCharacters).Methods("GET")

	l.Fatal(http.ListenAndServe(":8080", router))
}

func commonHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
