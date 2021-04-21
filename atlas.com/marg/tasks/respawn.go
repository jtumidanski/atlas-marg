package tasks

import (
	"atlas-marg/map"
	"atlas-marg/models"
	"atlas-marg/processor"
	"atlas-marg/registries"
	"github.com/sirupsen/logrus"
	"math"
	"math/rand"
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
	mks := registries.GetMapCharacterRegistry().GetMapsWithCharacters()
	for _, mk := range mks {
		go func(l logrus.FieldLogger, mk models.MapKey) {
			respawn(l, mk.WorldId, mk.ChannelId, mk.MapId)
		}(r.l, mk)
	}
}

func respawn(l logrus.FieldLogger, worldId byte, channelId byte, mapId int) {
	c := len(registries.GetMapCharacterRegistry().GetCharactersInMap(worldId, channelId, mapId))
	if c > 0 {
		sps, err := _map.NewMap(l).GetMonsterSpawnPoints(mapId)
		if err != nil {
			l.WithError(err).Errorf("Unable to get spawn points for map %d.", mapId)
			return
		}

		var ableSps []models.MonsterSpawnPoint
		for _, x := range sps {
			if x.MobTime >= 0 {
				ableSps = append(ableSps, x)
			}
		}

		monstersInMap, err := processor.NewMonster(l).CountInMap(worldId, channelId, mapId)
		if err != nil {
			l.WithError(err).Warnf("Assuming no monsters in map.")
		}

		monstersMax := getMonsterMax(c, len(ableSps))

		toSpawn := monstersMax - monstersInMap
		if toSpawn > 0 {
			result := shuffle(ableSps)
			for i := 0; i < toSpawn; i++ {
				x := result[i]
				processor.NewMonster(l).CreateMonster(worldId, channelId, mapId, x.Id, x.X, x.Y, x.Fh, x.Team)
			}
		}
	}
}

func shuffle(vals []models.MonsterSpawnPoint) []models.MonsterSpawnPoint {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	ret := make([]models.MonsterSpawnPoint, len(vals))
	perm := r.Perm(len(vals))
	for i, randIndex := range perm {
		ret[i] = vals[randIndex]
	}
	return ret
}

func getMonsterMax(characterCount int, spawnPointCount int) int {
	spawnRate := 0.70 + (0.05 * math.Min(6, float64(characterCount)))
	return int(math.Ceil(spawnRate * float64(spawnPointCount)))
}

func (r *Respawn) SleepTime() time.Duration {
	return time.Millisecond * time.Duration(r.interval)
}
