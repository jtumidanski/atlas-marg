package _map

import (
	"atlas-marg/models"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

func GetMonsterSpawnPoints(l logrus.FieldLogger, span opentracing.Span) func(mapId uint32) ([]models.MonsterSpawnPoint, error) {
	return func(mapId uint32) ([]models.MonsterSpawnPoint, error) {
		data, err := requestMonsterSpawnPoints(l, span)(mapId)
		if err != nil {
			return nil, err
		}
		return makeModel(data), nil
	}
}

func makeModel(data *MonsterInformationListDataContainer) []models.MonsterSpawnPoint {
	var result []models.MonsterSpawnPoint
	for _, x := range data.Data {
		result = append(result, models.MonsterSpawnPoint{
			Id:      x.Attributes.Id,
			MobTime: x.Attributes.MobTime,
			Team:    x.Attributes.Team,
			Cy:      x.Attributes.CY,
			F:       x.Attributes.F,
			Fh:      x.Attributes.FH,
			Rx0:     x.Attributes.RX0,
			Rx1:     x.Attributes.RX1,
			X:       x.Attributes.X,
			Y:       x.Attributes.Y,
		})
	}

	return result
}
