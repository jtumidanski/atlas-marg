package character

import (
	"atlas-marg/kafka"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type mapCharacterEvent struct {
	WorldId     byte   `json:"worldId"`
	ChannelId   byte   `json:"channelId"`
	MapId       uint32 `json:"mapId"`
	CharacterId uint32 `json:"characterId"`
	Type        string `json:"type"`
}

func emitEnterMap(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, characterId uint32) {
	producer := kafka.ProduceEvent(l, span, "TOPIC_MAP_CHARACTER_EVENT")
	return func(worldId byte, channelId byte, mapId uint32, characterId uint32) {
		e := &mapCharacterEvent{WorldId: worldId, ChannelId: channelId, MapId: mapId, CharacterId: characterId, Type: "ENTER"}
		producer(kafka.CreateKey(int(mapId)), e)
	}
}

func emitExitMap(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, characterId uint32) {
	producer := kafka.ProduceEvent(l, span, "TOPIC_MAP_CHARACTER_EVENT")
	return func(worldId byte, channelId byte, mapId uint32, characterId uint32) {
		e := &mapCharacterEvent{WorldId: worldId, ChannelId: channelId, MapId: mapId, CharacterId: characterId, Type: "EXIT"}
		producer(kafka.CreateKey(int(mapId)), e)
	}
}
