package processor

import (
	"atlas-marg/rest/attributes"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
)

func CountInMap(l logrus.FieldLogger) func(worldId byte, channelId byte, mapId int) (int, error) {
	return func(worldId byte, channelId byte, mapId int) (int, error) {
		r, err := http.Get(fmt.Sprintf("http://atlas-nginx:80/ms/morg/worlds/%d/channels/%d/maps/%d/monsters", worldId, channelId, mapId))
		if err != nil {
			l.WithError(err).Errorf("Retrieving monster data for map %d", mapId)
			return 0, err
		}

		td := &attributes.MonsterListDataContainer{}
		err = attributes.FromJSON(td, r.Body)
		if err != nil {
			l.WithError(err).Errorf("Decoding monster data for map %d", mapId)
			return 0, err
		}
		return len(td.Data), nil
	}
}

func CreateMonster(l logrus.FieldLogger) func(worldId byte, channelId byte, mapId int, monsterId int, x int, y int, fh int, team int) {
	return func(worldId byte, channelId byte, mapId int, monsterId int, x int, y int, fh int, team int) {
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
			l.WithError(err).Errorf("Marshalling monster")
		}

		_, err = http.Post(fmt.Sprintf("http://atlas-nginx:80/ms/morg/worlds/%d/channels/%d/maps/%d/monsters", worldId, channelId, mapId),
			"application/json; charset=utf-8", bytes.NewBuffer(jsonReq))
		if err != nil {
			l.WithError(err).Errorf("Creating monster for map %d", mapId)
		}
	}
}
