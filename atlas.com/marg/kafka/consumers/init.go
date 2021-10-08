package consumers

import (
	"atlas-marg/kafka/handler"
	"context"
	"github.com/sirupsen/logrus"
	"sync"
)

const (
	CharacterStatusEvent = "character_status_event"
	ChangeMapEvent       = "change_map_event"
)

func CreateEventConsumers(l *logrus.Logger, ctx context.Context, wg *sync.WaitGroup) {
	cec := func(topicToken string, name string, emptyEventCreator handler.EmptyEventCreator, processor handler.EventHandler) {
		createEventConsumer(l, ctx, wg, name, topicToken, emptyEventCreator, processor)
	}
	cec("TOPIC_CHARACTER_STATUS", CharacterStatusEvent, CharacterStatusCreator(), HandleCharacterStatus())
	cec("TOPIC_CHANGE_MAP_EVENT", ChangeMapEvent, MapChangedCreator(), HandleMapChanged())
}

func createEventConsumer(l *logrus.Logger, ctx context.Context, wg *sync.WaitGroup, name string, topicToken string, emptyEventCreator handler.EmptyEventCreator, processor handler.EventHandler) {
	wg.Add(1)
	go NewConsumer(l, ctx, wg, name, topicToken, "Map Registry Service", emptyEventCreator, processor)
}
