package cashshop

import (
	"atlas-marg/kafka"
	_map "atlas-marg/map"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

const (
	consumerNameEntryPoll = "cash_shop_entry_poll"
	topicTokenEntryPoll   = "TOPIC_CASH_SHOP_ENTRY_POLL"
)

func EntryPollConsumer(serviceName string) func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[cashShopEntryPoll](consumerNameEntryPoll, topicTokenEntryPoll, groupId, handleEntryPoll(serviceName))
	}
}

type cashShopEntryPoll struct {
	WorldId     byte   `json:"world_id"`
	ChannelId   byte   `json:"channel_id"`
	CharacterId uint32 `json:"character_id"`
}

func handleEntryPoll(serviceName string) kafka.HandlerFunc[cashShopEntryPoll] {
	return func(l logrus.FieldLogger, span opentracing.Span, command cashShopEntryPoll) {
		in, err := _map.InMiniDungeon(l, span)(command.CharacterId)
		if in || err != nil {
			RejectEntryPoll(l, span)(command.CharacterId, serviceName)
			return
		}
		ApproveEntryPoll(l, span)(command.CharacterId, serviceName)
	}
}
