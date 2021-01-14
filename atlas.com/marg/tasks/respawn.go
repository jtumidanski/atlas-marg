package tasks

import (
	"atlas-marg/models"
	"atlas-marg/processor"
	"atlas-marg/registries"
	"log"
	"math"
	"math/rand"
	"time"
)

type Respawn struct {
	l        *log.Logger
	interval int
}

func NewRespawn(l *log.Logger, interval int) *Respawn {
	return &Respawn{l, interval}
}

func (r *Respawn) Run() {
	mks := registries.GetMapCharacterRegistry().GetMapsWithCharacters()
	for _, mk := range mks {
		go func(l *log.Logger, mk models.MapKey) {
			respawn(l, mk.WorldId, mk.ChannelId, mk.MapId)
		}(r.l, mk)
	}
}

func respawn(l *log.Logger, worldId byte, channelId byte, mapId int) {
	c := len(registries.GetMapCharacterRegistry().GetCharactersInMap(worldId, channelId, mapId))
	if c > 0 {
		sps, err := processor.NewMap(l).GetMonsterSpawnPoints(mapId)
		if err == nil {
			var ableSps []models.MonsterSpawnPoint
			for _, x := range sps {
				if x.MobTime >= 0 {
					ableSps = append(ableSps, x)
				}
			}

			monstersInMap, err := processor.NewMonster(l).CountInMap(worldId, channelId, mapId)
			if err != nil {
				l.Print("Assuming no monsters in map.")
			}

			monstersMax := getMonsterMax(c, len(ableSps))
			if monstersMax-monstersInMap > 0 {
				rand.Seed(time.Now().UnixNano())

				dest := make([]models.MonsterSpawnPoint, len(ableSps))
				perm := rand.Perm(len(ableSps))
				for i, v := range perm {
					//goland:noinspection GoNilness
					dest[v] = ableSps[i]
				}

				for _, x := range dest {
					processor.NewMonster(l).CreateMonster(worldId, channelId, mapId, x.Id, x.X, x.Y, x.Fh, x.Team)
				}
			}
		}
	}
}

func getMonsterMax(characterCount int, spawnPointCount int) int {
	spawnRate := 0.70 + (0.05 * math.Min(6, float64(characterCount)))
	return int(math.Ceil(spawnRate * float64(spawnPointCount)))
}

func (r *Respawn) SleepTime() time.Duration {
	return time.Millisecond * time.Duration(r.interval)
}
