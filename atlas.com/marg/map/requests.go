package _map

import (
	"atlas-marg/rest/requests"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

const (
	mapInformationServicePrefix string = "/ms/mis/"
	mapInformationService              = requests.BaseRequest + mapInformationServicePrefix
	mapsResource                       = mapInformationService + "maps/"
	monstersResource                   = mapsResource + "%d/monsters"
)

func requestMonsterSpawnPoints(l logrus.FieldLogger, span opentracing.Span) func(mapId uint32) (*MonsterInformationListDataContainer, error) {
	return func(mapId uint32) (*MonsterInformationListDataContainer, error) {
		td := &MonsterInformationListDataContainer{}
		err := requests.Get(l, span)(fmt.Sprintf(monstersResource, mapId), td)
		if err != nil {
			l.WithError(err).Errorf("Retrieving monster spawn data for map %d", mapId)
			return nil, err
		}
		return td, nil
	}
}
