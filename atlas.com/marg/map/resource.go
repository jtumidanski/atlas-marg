package _map

import (
	"atlas-marg/json"
	"atlas-marg/map/character"
	"atlas-marg/rest"
	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

const GetMapCharacters = "get_map_characters"

func InitResource(router *mux.Router, l logrus.FieldLogger) {
	mRouter := router.PathPrefix("/worlds").Subrouter()
	mRouter.HandleFunc("/{worldId}/channels/{channelId}/maps/{mapId}/characters", registerGetMapCharacters(l)).Methods(http.MethodGet)
}

func registerGetMapCharacters(l logrus.FieldLogger) http.HandlerFunc {
	return rest.RetrieveSpan(GetMapCharacters, func(span opentracing.Span) http.HandlerFunc {
		return parseMap(l, func(worldId byte, channelId byte, mapId uint32) http.HandlerFunc {
			return handleGetMapCharacters(l)(span)(worldId, channelId, mapId)
		})
	})
}

type mapHandler func(worldId byte, channelId byte, mapId uint32) http.HandlerFunc

func parseMap(l logrus.FieldLogger, next mapHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		worldId, err := strconv.ParseUint(vars["worldId"], 10, 8)
		if err != nil {
			l.WithError(err).Errorf("Error parsing worldId as byte")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		channelId, err := strconv.ParseUint(vars["channelId"], 10, 8)
		if err != nil {
			l.WithError(err).Errorf("Error parsing channelId as byte")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		mapId, err := strconv.ParseUint(vars["mapId"], 10, 32)
		if err != nil {
			l.WithError(err).Errorf("Error parsing mapId as uint32")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		next(byte(worldId), byte(channelId), uint32(mapId))(w, r)
	}
}

func handleGetMapCharacters(l logrus.FieldLogger) func(span opentracing.Span) func(worldId byte, channelId byte, mapId uint32) http.HandlerFunc {
	return func(span opentracing.Span) func(worldId byte, channelId byte, mapId uint32) http.HandlerFunc {
		return func(worldId byte, channelId byte, mapId uint32) http.HandlerFunc {
			return func(rw http.ResponseWriter, r *http.Request) {
				var response CharacterDataListContainer
				response.Data = make([]CharactersDataBody, 0)

				for _, x := range character.GetRegistry().GetInMap(worldId, channelId, mapId) {
					var serverData = getMapCharactersResponseObject(x)
					response.Data = append(response.Data, serverData)
				}

				err := json.ToJSON(response, rw)
				if err != nil {
					l.WithError(err).Errorf("Error encoding GetChannelServers response")
					rw.WriteHeader(http.StatusInternalServerError)
				}
			}
		}
	}
}

func getMapCharactersResponseObject(id uint32) CharactersDataBody {
	return CharactersDataBody{
		Id:         strconv.Itoa(int(id)),
		Type:       "com.atlas.mrg.rest.attribute.MapCharacterAttributes",
		Attributes: CharacterAttributes{},
	}
}
