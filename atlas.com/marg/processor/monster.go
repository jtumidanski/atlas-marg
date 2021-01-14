package processor

import (
	"atlas-marg/attributes"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Monster struct {
	l *log.Logger
}

func NewMonster(l *log.Logger) *Monster {
	return &Monster{l}
}

func (m *Monster) CountInMap(worldId byte, channelId byte, mapId int) (int, error) {
	r, err := http.Get(fmt.Sprintf("http://atlas-nginx:80/ms/morg/worlds/%d/channels/%d/maps/%d/monsters", worldId, channelId, mapId))
	if err != nil {
		m.l.Printf("[ERROR] retrieving monster data for map %d", mapId)
		return 0, err
	}

	td := &attributes.MonsterListDataContainer{}
	err = attributes.FromJSON(td, r.Body)
	if err != nil {
		m.l.Printf("[ERROR] decoding monster data for map %d", mapId)
		return 0, err
	}
	return len(td.Data), nil
}

func (m *Monster) CreateMonster(worldId byte, channelId byte, mapId int, monsterId int, x int, y int, fh int, team int) {
	monster := attributes.MonsterInputDataContainer{
		Data: attributes.MonsterData{
			Id:   "0",
			Type: "com.atlas.morg.rest.attribute.MonsterAttributes",
			Attributes: attributes.MonsterAttributes{
				MonsterId: monsterId,
				X:         x,
				Y:         y,
				Team:      team,
				Fh:        fh,
			},
		},
	}
	jsonReq, err := json.Marshal(monster)
	if err != nil {
		m.l.Fatal("[ERROR] marshalling monster")
	}

	_, err = http.Post(fmt.Sprintf("http://atlas-nginx:80/ms/morg/worlds/%d/channels/%d/maps/%d/monsters", worldId, channelId, mapId),
		"application/json; charset=utf-8", bytes.NewBuffer(jsonReq))
	if err != nil {
		m.l.Printf("[ERROR] creating monster for map %d", mapId)
	}
}
