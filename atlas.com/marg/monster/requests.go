package monster

import (
	"atlas-marg/rest/requests"
	"fmt"
)

const (
	monsterRegistryServicePrefix string = "/ms/morg/"
	monsterRegistryService              = requests.BaseRequest + monsterRegistryServicePrefix
	mapMonstersResource                 = monsterRegistryService + "worlds/%d/channels/%d/maps/%d/monsters"
)

func requestInMap(worldId byte, channelId byte, mapId uint32) requests.Request[Attributes] {
	return requests.MakeGetRequest[Attributes](fmt.Sprintf(mapMonstersResource, worldId, channelId, mapId))
}

func requestCreate(worldId byte, channelId byte, mapId uint32, monsterId uint32, x int16, y int16, fh uint16, team int32) requests.PostRequest[Attributes] {
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
	return requests.MakePostRequest[Attributes](fmt.Sprintf(mapMonstersResource, worldId, channelId, mapId), monster)
}
