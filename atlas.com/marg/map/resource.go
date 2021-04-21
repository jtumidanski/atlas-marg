package _map

import (
	"atlas-marg/registries"
	"atlas-marg/rest/attributes"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func GetMapCharacters(l *log.Logger) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var response attributes.MapCharactersListDataContainer
		response.Data = make([]attributes.MapCharactersData, 0)

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
		mapId := value

		for _, x := range registries.GetMapCharacterRegistry().GetCharactersInMap(worldId, channelId, mapId) {
			var serverData = getMapCharactersResponseObject(x)
			response.Data = append(response.Data, serverData)
		}

		err = attributes.ToJSON(response, rw)
		if err != nil {
			l.WithError(err).Errorf("Error encoding GetChannelServers response")
			rw.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func getMapCharactersResponseObject(id int) attributes.MapCharactersData {
	return attributes.MapCharactersData{
		Id:         strconv.Itoa(id),
		Type:       "com.atlas.mrg.rest.attribute.MapCharacterAttributes",
		Attributes: attributes.MapCharactersAttributes{},
	}
}
