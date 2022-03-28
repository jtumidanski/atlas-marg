package _map

import (
	"atlas-marg/rest/requests"
	"fmt"
)

const (
	mapInformationServicePrefix string = "/ms/mis/"
	mapInformationService              = requests.BaseRequest + mapInformationServicePrefix
	mapsResource                       = mapInformationService + "maps/"
	monstersResource                   = mapsResource + "%d/monsters"
)

func requestMonsterSpawnPoints(mapId uint32) requests.Request[monsterAttributes] {
	return requests.MakeGetRequest[monsterAttributes](fmt.Sprintf(monstersResource, mapId))
}
