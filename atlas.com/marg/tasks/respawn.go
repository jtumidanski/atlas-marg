package tasks

import (
	"atlas-marg/map"
	"atlas-marg/map/monster"
	"github.com/sirupsen/logrus"
	"time"
)

type Respawn struct {
	l        logrus.FieldLogger
	interval int
}

func NewRespawn(l logrus.FieldLogger, interval int) *Respawn {
	return &Respawn{l, interval}
}

func (r *Respawn) Run() {
	mks := _map.GetCharacterRegistry().GetMapsWithCharacters()
	for _, mk := range mks {
		go monster.Spawn(r.l)(mk.WorldId, mk.ChannelId, mk.MapId)
	}
}

func (r *Respawn) SleepTime() time.Duration {
	return time.Millisecond * time.Duration(r.interval)
}
