// Package classification of Map API
//
// Documentation for Map API
//
// Schemes: http
// BasePath: /ms/mrg/worlds
// Version: 1.0.0
//
// Consumes:
// -application/json
//
// Produces:
// -application/json
// swagger:meta
package handlers

import (
	"atlas-marg/attributes"
	"atlas-marg/registries"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

// Map handler for map related queries
type Map struct {
	l *log.Logger
}

func NewMap(l *log.Logger) *Map {
	return &Map{l}
}

// swagger:route GET /worlds/{worldId}/channels/{channelId}/maps/{mapId}/characters maps getMapCharacters
// Return a list of characters in the map
// responses:
//	200: charactersResponse

// GetMapCharacters handles GET requests
func (m *Map) GetMapCharacters(rw http.ResponseWriter, r *http.Request) {

	var response attributes.MapCharactersListDataContainer
	response.Data = make([]attributes.MapCharactersData, 0)

	for _, x := range registries.GetMapCharacterRegistry().GetCharactersInMap(getWorldId(r), getChannelId(r), getMapId(r)) {
		var serverData = getMapCharactersResponseObject(x)
		response.Data = append(response.Data, serverData)
	}

	err := attributes.ToJSON(response, rw)
	if err != nil {
		m.l.Println("Error encoding GetChannelServers response")
		rw.WriteHeader(http.StatusInternalServerError)
	}
}

func getMapCharactersResponseObject(id int) attributes.MapCharactersData {
	return attributes.MapCharactersData{
		Id:         strconv.Itoa(id),
		Type:       "com.atlas.mrg.rest.attribute.MapCharacterAttributes",
		Attributes: attributes.MapCharactersAttributes{},
	}
}

func getWorldId(r *http.Request) byte {
	vars := mux.Vars(r)
	value, err := strconv.Atoi(vars["worldId"])
	if err != nil {
		log.Println("Error parsing worldId as integer")
		return 0
	}
	return byte(value)
}

func getChannelId(r *http.Request) byte {
	vars := mux.Vars(r)
	value, err := strconv.Atoi(vars["channelId"])
	if err != nil {
		log.Println("Error parsing channelId as integer")
		return 0
	}
	return byte(value)
}

func getMapId(r *http.Request) int {
	vars := mux.Vars(r)
	value, err := strconv.Atoi(vars["mapId"])
	if err != nil {
		log.Println("Error parsing mapId as integer")
		return 0
	}
	return value
}
