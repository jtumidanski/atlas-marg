package character

import (
	"atlas-marg/kafka"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

const (
	consumerNameStatus    = "character_status_event"
	consumerNameChangeMap = "change_map_event"
	topicNameStatus       = "TOPIC_CHARACTER_STATUS"
	topicNameChangeMap    = "TOPIC_CHANGE_MAP_EVENT"
)

func StatusConsumer(groupId string) kafka.ConsumerConfig {
	return kafka.NewConsumerConfig[statusEvent](consumerNameStatus, topicNameStatus, groupId, handleStatus())
}

type statusEvent struct {
	WorldId     byte   `json:"worldId"`
	ChannelId   byte   `json:"channelId"`
	AccountId   uint32 `json:"accountId"`
	CharacterId uint32 `json:"characterId"`
	Type        string `json:"type"`
}

func handleStatus() kafka.HandlerFunc[statusEvent] {
	return func(l logrus.FieldLogger, span opentracing.Span, event statusEvent) {
		if event.Type == "LOGIN" {
			EnterMap(l, span)(event.WorldId, event.ChannelId, event.CharacterId)
		} else if event.Type == "LOGOUT" {
			ExitMap(l, span)(event.WorldId, event.ChannelId, event.CharacterId)
		} else {
			l.Errorf("Unhandled event status %s.", event.Type)
		}
	}
}

func MapChangedConsumer(groupId string) kafka.ConsumerConfig {
	return kafka.NewConsumerConfig[mapChangedEvent](consumerNameChangeMap, topicNameChangeMap, groupId, HandleMapChanged())
}

type mapChangedEvent struct {
	WorldId     byte   `json:"worldId"`
	ChannelId   byte   `json:"channelId"`
	MapId       uint32 `json:"mapId"`
	PortalId    uint32 `json:"portalId"`
	CharacterId uint32 `json:"characterId"`
}

func HandleMapChanged() kafka.HandlerFunc[mapChangedEvent] {
	return func(l logrus.FieldLogger, span opentracing.Span, event mapChangedEvent) {
		TransitionMap(l, span)(event.WorldId, event.ChannelId, event.MapId, event.CharacterId)
	}
}
