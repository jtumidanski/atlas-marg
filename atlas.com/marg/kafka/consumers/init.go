package consumers

import (
	"context"
	"github.com/sirupsen/logrus"
)

func CreateEventConsumers(l *logrus.Logger) {
	cec := func(topicToken string, emptyEventCreator EmptyEventCreator, processor EventProcessor) {
		createEventConsumer(l, topicToken, emptyEventCreator, processor)
	}
	cec("TOPIC_CHARACTER_STATUS", CharacterStatusCreator(), HandleCharacterStatus())
	cec("TOPIC_CHANGE_MAP_EVENT", MapChangedCreator(), HandleMapChanged())
}

func createEventConsumer(l *logrus.Logger, topicToken string, emptyEventCreator EmptyEventCreator, processor EventProcessor) {
	h := func(logger logrus.FieldLogger, event interface{}) {
		processor(logger, event)
	}

	c := NewConsumer(l, context.Background(), h,
		SetGroupId("Map Registry Service"),
		SetTopicToken(topicToken),
		SetEmptyEventCreator(emptyEventCreator))
	go c.Init()
}
