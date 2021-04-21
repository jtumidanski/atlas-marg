package _map

import (
	"atlas-marg/models"
	"atlas-marg/rest/attributes"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Map struct {
	l logrus.FieldLogger
}

func NewMap(l logrus.FieldLogger) *Map {
	return &Map{l}
}

func (c *Map) GetMonsterSpawnPoints(mapId int) ([]models.MonsterSpawnPoint, error) {
	r, err := http.Get(fmt.Sprintf("http://atlas-nginx:80/ms/mis/maps/%d/monsters", mapId))
	if err != nil {
		c.l.WithError(err).Errorf("Retrieving monster spawn data for map %d", mapId)
		return nil, err
	}

	td := &attributes.MonsterInformationListDataContainer{}
	err = attributes.FromJSON(td, r.Body)
	if err != nil {
		c.l.WithError(err).Errorf("Decoding monster spawn data for map %d", mapId)
		return nil, err
	}

	var result []models.MonsterSpawnPoint
	for _, x := range td.Data {
		result = append(result, models.MonsterSpawnPoint{
			Id:      x.Attributes.MonsterId,
			MobTime: x.Attributes.MobTime,
			Team:    x.Attributes.Team,
			Cy:      x.Attributes.Cy,
			F:       x.Attributes.F,
			Fh:      x.Attributes.Fh,
			Rx0:     x.Attributes.Rx0,
			Rx1:     x.Attributes.Rx1,
			X:       x.Attributes.X,
			Y:       x.Attributes.Y,
		})
	}

	return result, nil
}
