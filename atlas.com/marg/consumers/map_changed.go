package consumers

import (
	"atlas-marg/events"
	"atlas-marg/processor"
	"atlas-marg/producers"
	"atlas-marg/registries"
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
	"time"
)

type MapChanged struct {
	l   *log.Logger
	ctx context.Context
}

func NewMapChanged(l *log.Logger, ctx context.Context) *MapChanged {
	return &MapChanged{l, ctx}
}

func (mc *MapChanged) Init() {
	mc.l.Printf("[INFO] [MapChanged] consumer started")

	t := processor.NewTopic(mc.l)
	td, err := t.GetTopic("TOPIC_CHANGE_MAP_EVENT")
	if err != nil {
		mc.l.Fatal("[ERROR] Unable to retrieve topic for consumer.")
	}

	mc.l.Printf("[INFO] [MapChanged] topic configuration retrieved")

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{os.Getenv("BOOTSTRAP_SERVERS")},
		Topic:   td.Attributes.Name,
		GroupID: "Map Registry",
		MaxWait: 50 * time.Millisecond,
	})

	mc.l.Printf("[INFO] [MapChanged] reading loop started")

	for {
		msg, err := r.ReadMessage(mc.ctx)
		if err != nil {
			panic("Could not successfully read message " + err.Error())
		}

		mc.l.Printf("[INFO] Handling [MapChanged] %s", msg.Value)

		var event events.MapChangedEvent
		err = json.Unmarshal(msg.Value, &event)
		if err != nil {
			mc.l.Println("Could not unmarshal event into event class ", msg.Value)
		} else {
			mc.processEvent(event)
		}

		mc.l.Printf("[INFO] Completed handling of [MapChanged] %s", msg.Value)
	}
}

func (mc *MapChanged) processEvent(event events.MapChangedEvent) {
	r := registries.GetMapCharacterRegistry()

	mk, err := r.GetMapForCharacter(event.CharacterId)
	if err == nil {
		r.RemoveCharacterFromMap(event.CharacterId)
		producers.NewMapCharacter(mc.l, context.Background()).EmitExit(event.WorldId, event.ChannelId, mk, event.CharacterId)
	}
	r.AddCharacterToMap(event.WorldId, event.ChannelId, event.MapId, event.CharacterId)

	mk, err = r.GetMapForCharacter(event.CharacterId)
	if err == nil {
		producers.NewMapCharacter(mc.l, context.Background()).EmitEnter(event.WorldId, event.ChannelId, mk, event.CharacterId)
	}
}
