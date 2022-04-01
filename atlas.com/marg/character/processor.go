package character

import (
	"atlas-marg/map/character"
	"atlas-marg/rest/requests"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

func GetMapId(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32) (uint32, error) {
	return func(characterId uint32) (uint32, error) {
		return requests.Provider[attributes, uint32](l, span)(requestById(characterId), mapIdGetter)()
	}
}

func mapIdGetter(body requests.DataBody[attributes]) (uint32, error) {
	attr := body.Attributes
	return attr.MapId, nil
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
		character.GetRegistry().AddToMap(worldId, channelId, mapId, characterId)
		emitEnterMap(l, span)(worldId, channelId, mapId, characterId)
	}
}

func ExitMap(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, characterId uint32) {
	return func(worldId byte, channelId byte, characterId uint32) {
		mk, err := character.GetRegistry().GetMapId(characterId)
		if err == nil {
			character.GetRegistry().RemoveFromMap(characterId)
			emitExitMap(l, span)(worldId, channelId, mk, characterId)
		}
	}
}
