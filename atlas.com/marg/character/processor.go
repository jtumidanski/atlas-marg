package character

import (
	"atlas-marg/kafka/producers"
	"atlas-marg/map"
	"atlas-marg/map/monster"
	"atlas-marg/reactor"
	"github.com/sirupsen/logrus"
)

func GetMapId(l logrus.FieldLogger) func(characterId uint32) (uint32, error) {
	return func(characterId uint32) (uint32, error) {
		c, err := requestById(l)(characterId)
		if err != nil {
			l.WithError(err).Errorf("Unable to locate the character.")
			return 0, err
		}
		return c.Data.Attributes.MapId, nil
	}
}

func TransitionMap(l logrus.FieldLogger) func(worldId byte, channelId byte, mapId uint32, characterId uint32) {
	return func(worldId byte, channelId byte, mapId uint32, characterId uint32) {
		ExitMap(l)(worldId, channelId, characterId)
		enterMap(l)(worldId, channelId, mapId, characterId)
	}
}

func EnterMap(l logrus.FieldLogger) func(worldId byte, channelId byte, characterId uint32) {
	return func(worldId byte, channelId byte, characterId uint32) {
		mk, err := GetMapId(l)(characterId)
		if err != nil {
			l.WithError(err).Errorf("Error retrieving intial map for character %d.", characterId)
			return
		}
		enterMap(l)(worldId, channelId, mk, characterId)
	}
}

func enterMap(l logrus.FieldLogger) func(worldId byte, channelId byte, mapId uint32, characterId uint32) {
	return func(worldId byte, channelId byte, mapId uint32, characterId uint32) {
		_map.GetCharacterRegistry().AddToMap(worldId, channelId, mapId, characterId)
		monster.Spawn(l)(worldId, channelId, mapId)
		reactor.SpawnMissing(l)(worldId, channelId, mapId)
		producers.EnterMap(l)(worldId, channelId, mapId, characterId)
	}
}

func ExitMap(l logrus.FieldLogger) func(worldId byte, channelId byte, characterId uint32) {
	return func(worldId byte, channelId byte, characterId uint32) {
		mk, err := _map.GetCharacterRegistry().GetMapId(characterId)
		if err == nil {
			_map.GetCharacterRegistry().RemoveFromMap(characterId)
			producers.ExitMap(l)(worldId, channelId, mk, characterId)
		}
	}
}
