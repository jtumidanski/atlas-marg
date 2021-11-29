package tasks

import (
	"atlas-marg/map"
	"atlas-marg/map/monster"
	"atlas-marg/reactor"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"time"
)

const RespawnTask = "respawn_task"

type Respawn struct {
	l        logrus.FieldLogger
	interval int
}

func NewRespawn(l logrus.FieldLogger, interval int) *Respawn {
	return &Respawn{l, interval}
}

func (r *Respawn) Run() {
	span := opentracing.StartSpan(RespawnTask)
	mks := _map.GetCharacterRegistry().GetMapsWithCharacters()
	for _, mk := range mks {
		go monster.Spawn(r.l, span)(mk.WorldId, mk.ChannelId, mk.MapId)
		go reactor.SpawnMissing(r.l, span)(mk.WorldId, mk.ChannelId, mk.MapId)
	}
	span.Finish()
}

func (r *Respawn) SleepTime() time.Duration {
	return time.Millisecond * time.Duration(r.interval)
}
