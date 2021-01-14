package producers

import (
	"atlas-marg/events"
	"atlas-marg/processor"
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
	"strconv"
	"time"
)

type MapCharacter struct {
	l   *log.Logger
	ctx context.Context
}

func NewMapCharacter(l *log.Logger, ctx context.Context) *MapCharacter {
	return &MapCharacter{l, ctx}
}

func (m *MapCharacter) EmitEnter(worldId byte, channelId byte, mapId int, characterId int) {
	m.emit(worldId, channelId, mapId, characterId, "ENTER")
}

func (m *MapCharacter) EmitExit(worldId byte, channelId byte, mapId int, characterId int) {
	m.emit(worldId, channelId, mapId, characterId, "EXIT")
}

func (m *MapCharacter) emit(worldId byte, channelId byte, mapId int, characterId int, theType string) {
	t := processor.NewTopic(m.l)
	td, err := t.GetTopic("TOPIC_MAP_CHARACTER_EVENT")
	if err != nil {
		m.l.Fatal("[ERROR] Unable to retrieve topic for consumer.")
	}

	w := &kafka.Writer{
		Addr:         kafka.TCP(os.Getenv("BOOTSTRAP_SERVERS")),
		Topic:        td.Attributes.Name,
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 50 * time.Millisecond,
	}

	e := &events.MapCharacterEvent{WorldId: worldId, ChannelId: channelId, MapId: mapId, CharacterId: characterId, Type: theType}
	r, err := json.Marshal(e)
	if err != nil {
		m.l.Fatal("[ERROR] Unable to marshall event.")
	}

	m.l.Printf("[INFO] Sending [MapCharacterEvent] key %s", createKey(mapId))
	m.l.Printf("[INFO] Sending [MapCharacterEvent] value %s", r)

	err = w.WriteMessages(context.Background(), kafka.Message{
		Key:   createKey(mapId),
		Value: r,
	})
	if err != nil {
		m.l.Fatal("[ERROR] Unable to produce event.")
	}
}

func createKey(key int) []byte {
	var empty = make([]byte, 8)
	sk := []byte(strconv.Itoa(key))

	start := len(empty) - len(sk)
	var result = empty

	for i := 0; i < start; i++ {
		result[i] = 48
	}

	for i := start; i < len(empty); i++ {
		result[i] = sk[i-start]
	}
	return result
}
