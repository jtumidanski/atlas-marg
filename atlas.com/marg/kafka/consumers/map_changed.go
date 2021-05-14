package consumers

import (
	"atlas-marg/kafka/handler"
	"atlas-marg/kafka/producers"
	"atlas-marg/map/monster"
	"atlas-marg/registries"
	"github.com/sirupsen/logrus"
)

type mapChangedEvent struct {
	WorldId     byte `json:"worldId"`
	ChannelId   byte `json:"channelId"`
	MapId       int  `json:"mapId"`
	PortalId    int  `json:"portalId"`
	CharacterId int  `json:"characterId"`
}

func MapChangedCreator() handler.EmptyEventCreator {
	return func() interface{} {
		return &mapChangedEvent{}
	}
}

func HandleMapChanged() handler.EventHandler {
	return func(l logrus.FieldLogger, e interface{}) {
		if event, ok := e.(*mapChangedEvent); ok {
			r := registries.GetMapCharacterRegistry()

			mk, err := r.GetMapForCharacter(event.CharacterId)
			if err == nil {
				r.RemoveCharacterFromMap(event.CharacterId)
				producers.ExitMap(l)(event.WorldId, event.ChannelId, mk, event.CharacterId)
			}
			r.AddCharacterToMap(event.WorldId, event.ChannelId, event.MapId, event.CharacterId)

			mk, err = r.GetMapForCharacter(event.CharacterId)
			if err == nil {
				monster.Spawn(l)(event.WorldId, event.ChannelId, mk)
				producers.EnterMap(l)(event.WorldId, event.ChannelId, mk, event.CharacterId)
			}
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}
