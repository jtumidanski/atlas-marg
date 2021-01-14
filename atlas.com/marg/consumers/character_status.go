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

type CharacterStatus struct {
	l   *log.Logger
	ctx context.Context
}

func NewCharacterStatus(l *log.Logger, ctx context.Context) *CharacterStatus {
	return &CharacterStatus{l, ctx}
}

func (c *CharacterStatus) Init() {
	t := processor.NewTopic(c.l)
	td, err := t.GetTopic("TOPIC_CHARACTER_STATUS")
	if err != nil {
		c.l.Fatal("[ERROR] Unable to retrieve topic for consumer.")
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{os.Getenv("BOOTSTRAP_SERVERS")},
		Topic:   td.Attributes.Name,
		GroupID: "Map Registry",
		MaxWait: 50 * time.Millisecond,
	})
	for {
		msg, err := r.ReadMessage(c.ctx)
		if err != nil {
			panic("Could not successfully read message " + err.Error())
		}

		c.l.Printf("[INFO] Handling [CharacterStatusEvent] %s", msg.Value)

		var event events.CharacterStatusEvent
		err = json.Unmarshal(msg.Value, &event)
		if err != nil {
			c.l.Println("Could not unmarshal event into event class ", msg.Value)
		} else {
			c.processEvent(event)
		}
	}
}

func (c *CharacterStatus) processEvent(event events.CharacterStatusEvent) {
	if event.Type == "LOGIN" {
		mk, err := processor.NewCharacter(c.l).GetMapForCharacter(event.CharacterId)
		if err == nil {
			registries.GetMapCharacterRegistry().AddCharacterToMap(event.WorldId, event.ChannelId, mk, event.CharacterId)
			producers.NewMapCharacter(c.l, context.Background()).EmitEnter(event.WorldId, event.ChannelId, mk, event.CharacterId)
		}
	} else if event.Type == "LOGOUT" {
		mk, err := processor.NewCharacter(c.l).GetMapForCharacter(event.CharacterId)
		if err == nil {
			registries.GetMapCharacterRegistry().RemoveCharacterFromMap(event.CharacterId)
			producers.NewMapCharacter(c.l, context.Background()).EmitExit(event.WorldId, event.ChannelId, mk, event.CharacterId)
		}
	} else {
		c.l.Println("Unhandled event status ", event.Type)
	}
}
