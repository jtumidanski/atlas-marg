package monster

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

func requestSpawnPoints(mapId uint32) requests.Request[attributes] {
	return requests.MakeGetRequest[attributes](fmt.Sprintf(monstersResource, mapId))
}
