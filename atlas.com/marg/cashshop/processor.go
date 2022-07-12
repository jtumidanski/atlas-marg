package cashshop

import (
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

func RegisterGatekeeper(l logrus.FieldLogger, span opentracing.Span) func(service string) {
	return func(service string) {
		emitGatekeeperRegistration(l, span)(service)
	}
}

func UnregisterGatekeeper(l logrus.FieldLogger, span opentracing.Span) func(service string) {
	return func(service string) {
		emitGatekeeperUnregistering(l, span)(service)
	}
}

func ApproveEntryPoll(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32, service string) {
	return func(characterId uint32, service string) {
		l.Debugf("Approving character %d entry into the cash shop.", characterId)
		emitEntryPollApproval(l, span)(characterId, service)
	}
}

func RejectEntryPoll(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32, service string) {
	return func(characterId uint32, service string) {
		l.Debugf("Rejecting character %d entry into the cash shop.")
		emitEntryPollRejection(l, span)(characterId, service, "PINK_TEXT", "Entering Cash Shop is disabled when inside a Mini-Dungeon.")
	}
}
