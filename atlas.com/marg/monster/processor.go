package monster

import (
	"github.com/sirupsen/logrus"
)

func CountInMap(l logrus.FieldLogger) func(worldId byte, channelId byte, mapId uint32) (int, error) {
	return func(worldId byte, channelId byte, mapId uint32) (int, error) {
		data, err := requestInMap(l)(worldId, channelId, mapId)
		if err != nil {
			return 0, err
		}
		return len(data.Data), nil
	}
}

func CreateMonster(l logrus.FieldLogger) func(worldId byte, channelId byte, mapId uint32, monsterId uint32, x int16, y int16, fh uint16, team int32) {
	return func(worldId byte, channelId byte, mapId uint32, monsterId uint32, x int16, y int16, fh uint16, team int32) {
		err := requestCreate(worldId, channelId, mapId, monsterId, x, y, fh, team)
		if err != nil {
			l.WithError(err).Errorf("Creating monster for map %d", mapId)
		}
	}
}
