package cashshop

import (
	"atlas-marg/kafka"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type gatekeeperCommand struct {
	Service string `json:"service"`
	Type    string `json:"type"`
}

func emitGatekeeperRegistration(l logrus.FieldLogger, span opentracing.Span) func(service string) {
	return func(service string) {
		emitGatekeeperCommand(l, span)(service, "REGISTER")
	}
}

func emitGatekeeperUnregistering(l logrus.FieldLogger, span opentracing.Span) func(service string) {
	return func(service string) {
		emitGatekeeperCommand(l, span)(service, "UNREGISTER")
	}
}

func emitGatekeeperCommand(l logrus.FieldLogger, span opentracing.Span) func(service string, theType string) {
	producer := kafka.ProduceEvent(l, span, "TOPIC_CASH_SHOP_GATEKEEPER_COMMAND")
	return func(service string, theType string) {
		e := &gatekeeperCommand{
			Service: service,
			Type:    theType,
		}
		producer(kafka.CreateKey(0), e)
	}
}

type entryPollResponse struct {
	CharacterId uint32 `json:"character_id"`
	Service     string `json:"service"`
	Type        string `json:"type"`
	MessageType string `json:"message_type"`
	Message     string `json:"message"`
}

func emitEntryPollApproval(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32, service string) {
	return func(characterId uint32, service string) {
		emitEntryPollResponse(l, span)(characterId, service, "APPROVE", "", "")
	}
}

func emitEntryPollRejection(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32, service string, messageType string, message string) {
	return func(characterId uint32, service string, messageType string, message string) {
		emitEntryPollResponse(l, span)(characterId, service, "DENY", messageType, message)
	}
}

func emitEntryPollResponse(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32, service string, theType string, messageType string, message string) {
	producer := kafka.ProduceEvent(l, span, "TOPIC_CASH_SHOP_ENTRY_POLL_RESPONSE")
	return func(characterId uint32, service string, theType string, messageType string, message string) {
		e := &entryPollResponse{
			CharacterId: characterId,
			Service:     service,
			Type:        theType,
			MessageType: messageType,
			Message:     message,
		}
		producer(kafka.CreateKey(int(characterId)), e)
	}
}
