package consumers

import (
	"atlas-marg/kafka/producers"
	"atlas-marg/processor"
	"atlas-marg/registries"
	"context"
	"github.com/sirupsen/logrus"
)

type characterStatusEvent struct {
	WorldId     byte   `json:"worldId"`
	ChannelId   byte   `json:"channelId"`
	AccountId   int    `json:"accountId"`
	CharacterId int    `json:"characterId"`
	Type        string `json:"type"`
}

func CharacterStatusCreator() EmptyEventCreator {
	return func() interface{} {
		return &characterStatusEvent{}
	}
}

func HandleCharacterStatus() EventProcessor {
	return func(l logrus.FieldLogger, e interface{}) {
		if event, ok := e.(*characterStatusEvent); ok {
			if event.Type == "LOGIN" {
				mk, err := processor.NewCharacter(l).GetMapForCharacter(event.CharacterId)
				if err == nil {
					registries.GetMapCharacterRegistry().AddCharacterToMap(event.WorldId, event.ChannelId, mk, event.CharacterId)
					producers.NewMapCharacter(l, context.Background()).EmitEnter(event.WorldId, event.ChannelId, mk, event.CharacterId)
				}
			} else if event.Type == "LOGOUT" {
				mk, err := processor.NewCharacter(l).GetMapForCharacter(event.CharacterId)
				if err == nil {
					registries.GetMapCharacterRegistry().RemoveCharacterFromMap(event.CharacterId)
					producers.NewMapCharacter(l, context.Background()).EmitExit(event.WorldId, event.ChannelId, mk, event.CharacterId)
				}
			} else {
				l.Errorf("Unhandled event status %s.", event.Type)
			}
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}
