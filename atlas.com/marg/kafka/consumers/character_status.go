package consumers

import (
	"atlas-marg/character"
	"atlas-marg/kafka/handler"
	"github.com/opentracing/opentracing-go"
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
	return func(l logrus.FieldLogger, span opentracing.Span, e interface{}) {
		if event, ok := e.(*characterStatusEvent); ok {
			if event.Type == "LOGIN" {
				character.EnterMap(l, span)(event.WorldId, event.ChannelId, event.CharacterId)
			} else if event.Type == "LOGOUT" {
				character.ExitMap(l, span)(event.WorldId, event.ChannelId, event.CharacterId)
			} else {
				l.Errorf("Unhandled event status %s.", event.Type)
			}
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}
