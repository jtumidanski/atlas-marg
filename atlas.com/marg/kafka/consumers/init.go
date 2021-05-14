package consumers

import (
	"atlas-marg/kafka/handler"
	"github.com/sirupsen/logrus"
)

func CreateEventConsumers(l *logrus.Logger) {
	cec := func(topicToken string, emptyEventCreator handler.EmptyEventCreator, processor handler.EventHandler) {
		createEventConsumer(l, topicToken, emptyEventCreator, processor)
	}
	cec("TOPIC_CHARACTER_STATUS", CharacterStatusCreator(), HandleCharacterStatus())
	cec("TOPIC_CHANGE_MAP_EVENT", MapChangedCreator(), HandleMapChanged())
}

func createEventConsumer(l *logrus.Logger, topicToken string, emptyEventCreator handler.EmptyEventCreator, processor handler.EventHandler) {
	go NewConsumer(l, topicToken, "Map Registry Service", emptyEventCreator, processor)
}
