package _map

import (
	"atlas-marg/models"
	"atlas-marg/rest/requests"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

func GetMonsterSpawnPoints(l logrus.FieldLogger, span opentracing.Span) func(mapId uint32) ([]models.MonsterSpawnPoint, error) {
	return func(mapId uint32) ([]models.MonsterSpawnPoint, error) {
		data, err := requestMonsterSpawnPoints(mapId)(l, span)
		if err != nil {
			return nil, err
		}
		var result []models.MonsterSpawnPoint
		for _, d := range data.DataList() {
			result = append(result, makeModel(d))
		}
		return result, nil
	}
}

func makeModel(data requests.DataBody[monsterAttributes]) models.MonsterSpawnPoint {
	attr := data.Attributes
	return models.MonsterSpawnPoint{
		Id:      attr.Id,
		MobTime: attr.MobTime,
		Team:    attr.Team,
		Cy:      attr.CY,
		F:       attr.F,
		Fh:      attr.FH,
		Rx0:     attr.RX0,
		Rx1:     attr.RX1,
		X:       attr.X,
		Y:       attr.Y,
	}
}
