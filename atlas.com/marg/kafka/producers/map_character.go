package producers

import (
	"github.com/sirupsen/logrus"
)

type mapCharacterEvent struct {
	WorldId     byte   `json:"worldId"`
	ChannelId   byte   `json:"channelId"`
	MapId       uint32 `json:"mapId"`
	CharacterId uint32 `json:"characterId"`
	Type        string `json:"type"`
}

func EnterMap(l logrus.FieldLogger) func(worldId byte, channelId byte, mapId uint32, characterId uint32) {
	producer := ProduceEvent(l, "TOPIC_MAP_CHARACTER_EVENT")
	return func(worldId byte, channelId byte, mapId uint32, characterId uint32) {
		e := &mapCharacterEvent{WorldId: worldId, ChannelId: channelId, MapId: mapId, CharacterId: characterId, Type: "ENTER"}
		producer(CreateKey(mapId), e)
	}
}

func ExitMap(l logrus.FieldLogger) func(worldId byte, channelId byte, mapId uint32, characterId uint32) {
	producer := ProduceEvent(l, "TOPIC_MAP_CHARACTER_EVENT")
	return func(worldId byte, channelId byte, mapId uint32, characterId uint32) {
		e := &mapCharacterEvent{WorldId: worldId, ChannelId: channelId, MapId: mapId, CharacterId: characterId, Type: "EXIT"}
		producer(CreateKey(mapId), e)
	}
}
