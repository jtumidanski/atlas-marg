package monster

import (
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

func CountInMap(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32) (int, error) {
	return func(worldId byte, channelId byte, mapId uint32) (int, error) {
		data, err := requestInMap(l, span)(worldId, channelId, mapId)
		if err != nil {
			return 0, err
		}
		return len(data.Data), nil
	}
}

func CreateMonster(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, monsterId uint32, x int16, y int16, fh uint16, team int32) {
	return func(worldId byte, channelId byte, mapId uint32, monsterId uint32, x int16, y int16, fh uint16, team int32) {
		err := requestCreate(l, span)(worldId, channelId, mapId, monsterId, x, y, fh, team)
		if err != nil {
			l.WithError(err).Errorf("Creating monster for map %d", mapId)
		}
	}
}
