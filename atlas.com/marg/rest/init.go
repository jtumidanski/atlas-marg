package rest

import (
	_map "atlas-marg/map"
	"context"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"sync"
)

func CreateRestService(l *logrus.Logger, ctx context.Context, wg *sync.WaitGroup) {
	go NewServer(l, ctx, wg, ProduceRoutes)
}

func ProduceRoutes(l logrus.FieldLogger) http.Handler {
		router := mux.NewRouter().StrictSlash(true).PathPrefix("/ms/mrg").Subrouter()
		router.Use(CommonHeader)

		mRouter := router.PathPrefix("/worlds").Subrouter()
		mRouter.HandleFunc("/{worldId}/channels/{channelId}/maps/{mapId}/characters", _map.GetMapCharacters(l)).Methods(http.MethodGet)
		return router
}
