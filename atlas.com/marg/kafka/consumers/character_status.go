package consumers

import (
	"atlas-marg/character"
	"atlas-marg/kafka/handler"
	"atlas-marg/kafka/producers"
	"atlas-marg/registries"
	"github.com/sirupsen/logrus"
)

type characterStatusEvent struct {
	WorldId     byte   `json:"worldId"`
	ChannelId   byte   `json:"channelId"`
	AccountId   uint32 `json:"accountId"`
	CharacterId uint32 `json:"characterId"`
	Type        string `json:"type"`
}

func CharacterStatusCreator() handler.EmptyEventCreator {
	return func() interface{} {
		return &characterStatusEvent{}
	}
}

func HandleCharacterStatus() handler.EventHandler {
	return func(l logrus.FieldLogger, e interface{}) {
		if event, ok := e.(*characterStatusEvent); ok {
			if event.Type == "LOGIN" {
				mk, err := character.GetMapForCharacter(l)(event.CharacterId)
				if err == nil {
					registries.GetMapCharacterRegistry().AddCharacterToMap(event.WorldId, event.ChannelId, mk, event.CharacterId)
					producers.EnterMap(l)(event.WorldId, event.ChannelId, mk, event.CharacterId)
				}
			} else if event.Type == "LOGOUT" {
				mk, err := character.GetMapForCharacter(l)(event.CharacterId)
				if err == nil {
					registries.GetMapCharacterRegistry().RemoveCharacterFromMap(event.CharacterId)
					producers.ExitMap(l)(event.WorldId, event.ChannelId, mk, event.CharacterId)
				}
			} else {
				l.Errorf("Unhandled event status %s.", event.Type)
			}
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}
