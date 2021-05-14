package producers

import (
	"github.com/sirupsen/logrus"
)

type mapCharacterEvent struct {
	WorldId     byte   `json:"worldId"`
	ChannelId   byte   `json:"channelId"`
	MapId       int    `json:"mapId"`
	CharacterId int    `json:"characterId"`
	Type        string `json:"type"`
}

func EnterMap(l logrus.FieldLogger) func(worldId byte, channelId byte, mapId int, characterId int) {
	producer := ProduceEvent(l, "TOPIC_MAP_CHARACTER_EVENT")
	return func(worldId byte, channelId byte, mapId int, characterId int) {
		e := &mapCharacterEvent{WorldId: worldId, ChannelId: channelId, MapId: mapId, CharacterId: characterId, Type: "ENTER"}
		producer(CreateKey(mapId), e)
	}
}

func ExitMap(l logrus.FieldLogger) func(worldId byte, channelId byte, mapId int, characterId int) {
	producer := ProduceEvent(l, "TOPIC_MAP_CHARACTER_EVENT")
	return func(worldId byte, channelId byte, mapId int, characterId int) {
		e := &mapCharacterEvent{WorldId: worldId, ChannelId: channelId, MapId: mapId, CharacterId: characterId, Type: "EXIT"}
		producer(CreateKey(mapId), e)
	}
}
