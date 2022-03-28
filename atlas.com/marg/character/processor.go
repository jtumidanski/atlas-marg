package character

import (
	"atlas-marg/map"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

func GetMapId(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32) (uint32, error) {
	return func(characterId uint32) (uint32, error) {
		c, err := requestById(characterId)(l, span)
		if err != nil {
			l.WithError(err).Errorf("Unable to locate the character.")
			return 0, err
		}
		attr := c.Data().Attributes
		return attr.MapId, nil
	}
}

func TransitionMap(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, characterId uint32) {
	return func(worldId byte, channelId byte, mapId uint32, characterId uint32) {
		ExitMap(l, span)(worldId, channelId, characterId)
		enterMap(l, span)(worldId, channelId, mapId, characterId)
	}
}

func EnterMap(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, characterId uint32) {
	return func(worldId byte, channelId byte, characterId uint32) {
		mk, err := GetMapId(l, span)(characterId)
		if err != nil {
			l.WithError(err).Errorf("Error retrieving intial map for character %d.", characterId)
			return
		}
		enterMap(l, span)(worldId, channelId, mk, characterId)
	}
}

func enterMap(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, characterId uint32) {
	return func(worldId byte, channelId byte, mapId uint32, characterId uint32) {
		_map.GetCharacterRegistry().AddToMap(worldId, channelId, mapId, characterId)
		emitEnterMap(l, span)(worldId, channelId, mapId, characterId)
	}
}

func ExitMap(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, characterId uint32) {
	return func(worldId byte, channelId byte, characterId uint32) {
		mk, err := _map.GetCharacterRegistry().GetMapId(characterId)
		if err == nil {
			_map.GetCharacterRegistry().RemoveFromMap(characterId)
			emitExitMap(l, span)(worldId, channelId, mk, characterId)
		}
	}
}
