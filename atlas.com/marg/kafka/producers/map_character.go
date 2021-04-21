package producers

import (
	"context"
	"github.com/sirupsen/logrus"
)

type mapCharacterEvent struct {
	WorldId     byte   `json:"worldId"`
	ChannelId   byte   `json:"channelId"`
	MapId       int    `json:"mapId"`
	CharacterId int    `json:"characterId"`
	Type        string `json:"type"`
}

type MapCharacter struct {
	l   logrus.FieldLogger
	ctx context.Context
}

func NewMapCharacter(l logrus.FieldLogger, ctx context.Context) *MapCharacter {
	return &MapCharacter{l, ctx}
}

func (m *MapCharacter) EmitEnter(worldId byte, channelId byte, mapId int, characterId int) {
	m.emit(worldId, channelId, mapId, characterId, "ENTER")
}

func (m *MapCharacter) EmitExit(worldId byte, channelId byte, mapId int, characterId int) {
	m.emit(worldId, channelId, mapId, characterId, "EXIT")
}

func (m *MapCharacter) emit(worldId byte, channelId byte, mapId int, characterId int, theType string) {
	e := &mapCharacterEvent{WorldId: worldId, ChannelId: channelId, MapId: mapId, CharacterId: characterId, Type: theType}
	produceEvent(m.l, "TOPIC_MAP_CHARACTER_EVENT", createKey(mapId), e)
}
