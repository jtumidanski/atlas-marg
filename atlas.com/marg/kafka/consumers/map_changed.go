package consumers

import (
	"atlas-marg/kafka/producers"
	"atlas-marg/registries"
	"context"
	"github.com/sirupsen/logrus"
)

type mapChangedEvent struct {
	WorldId     byte `json:"worldId"`
	ChannelId   byte `json:"channelId"`
	MapId       int  `json:"mapId"`
	PortalId    int  `json:"portalId"`
	CharacterId int  `json:"characterId"`
}

func MapChangedCreator() EmptyEventCreator {
	return func() interface{} {
		return &mapChangedEvent{}
	}
}

func HandleMapChanged() EventProcessor {
	return func(l logrus.FieldLogger, e interface{}) {
		if event, ok := e.(*mapChangedEvent); ok {
			r := registries.GetMapCharacterRegistry()

			mk, err := r.GetMapForCharacter(event.CharacterId)
			if err == nil {
				r.RemoveCharacterFromMap(event.CharacterId)
				producers.NewMapCharacter(l, context.Background()).EmitExit(event.WorldId, event.ChannelId, mk, event.CharacterId)
			}
			r.AddCharacterToMap(event.WorldId, event.ChannelId, event.MapId, event.CharacterId)

			mk, err = r.GetMapForCharacter(event.CharacterId)
			if err == nil {
				producers.NewMapCharacter(l, context.Background()).EmitEnter(event.WorldId, event.ChannelId, mk, event.CharacterId)
			}
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}
