package monster

import (
	"atlas-marg/rest/requests"
	"fmt"
	"github.com/sirupsen/logrus"
)

const (
	monsterRegistryServicePrefix string = "/ms/morg/"
	monsterRegistryService              = requests.BaseRequest + monsterRegistryServicePrefix
	mapMonstersResource                 = monsterRegistryService + "worlds/%d/channels/%d/maps/%d/monsters"
)

func requestInMap(l logrus.FieldLogger) func(worldId byte, channelId byte, mapId uint32) (*DataListContainer, error) {
	return func(worldId byte, channelId byte, mapId uint32) (*DataListContainer, error) {
		ar := &DataListContainer{}
		err := requests.Get(l)(fmt.Sprintf(mapMonstersResource, worldId, channelId, mapId), ar)
		if err != nil {
			return nil, err
		}
		return ar, nil
	}
}

func requestCreate(worldId byte, channelId byte, mapId uint32, monsterId uint32, x int16, y int16, fh uint16, team int32) error {
	monster := InputDataContainer{
		Data: DataBody{
			Id:   "0",
			Type: "com.atlas.morg.rest.attribute.MonsterAttributes",
			Attributes: Attributes{
				MonsterId: monsterId,
				X:         x,
				Y:         y,
				Team:      team,
				Fh:        fh,
			},
		},
	}
	_, err := requests.Post(fmt.Sprintf(mapMonstersResource, worldId, channelId, mapId), monster)
	return err
}
