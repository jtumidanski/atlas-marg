package _map

import (
	"atlas-marg/json"
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
	mRouter.HandleFunc("/{worldId}/channels/{channelId}/maps/{mapId}/characters", rest.RetrieveSpan(GetMapCharacters, HandleGetMapCharacters(l))).Methods(http.MethodGet)
}

func HandleGetMapCharacters(l logrus.FieldLogger) rest.SpanHandler {
	return func(span opentracing.Span) http.HandlerFunc {
		return func(rw http.ResponseWriter, r *http.Request) {
			var response CharacterDataListContainer
			response.Data = make([]CharactersDataBody, 0)

			vars := mux.Vars(r)
			value, err := strconv.Atoi(vars["worldId"])
			if err != nil {
				l.WithError(err).Errorf("Error parsing worldId as integer")
				rw.WriteHeader(http.StatusBadRequest)
				return
			}
			worldId := byte(value)

			vars = mux.Vars(r)
			value, err = strconv.Atoi(vars["channelId"])
			if err != nil {
				l.WithError(err).Errorf("Error parsing channelId as integer")
				rw.WriteHeader(http.StatusBadRequest)
				return
			}
			channelId := byte(value)

			vars = mux.Vars(r)
			value, err = strconv.Atoi(vars["mapId"])
			if err != nil {
				l.WithError(err).Errorf("Error parsing mapId as integer")
				rw.WriteHeader(http.StatusBadRequest)
				return
			}
			mapId := uint32(value)

			for _, x := range GetCharacterRegistry().GetInMap(worldId, channelId, mapId) {
				var serverData = getMapCharactersResponseObject(x)
				response.Data = append(response.Data, serverData)
			}

			err = json.ToJSON(response, rw)
			if err != nil {
				l.WithError(err).Errorf("Error encoding GetChannelServers response")
				rw.WriteHeader(http.StatusInternalServerError)
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
