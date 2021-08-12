package monster

import (
	_map "atlas-marg/map"
	"atlas-marg/models"
	"atlas-marg/processor"
	"atlas-marg/registries"
	"github.com/sirupsen/logrus"
	"math"
	"math/rand"
	"time"
)

func Spawn(l logrus.FieldLogger) func(worldId byte, channelId byte, mapId uint32) {
	return func(worldId byte, channelId byte, mapId uint32) {
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

			monstersInMap, err := processor.CountInMap(l)(worldId, channelId, mapId)
			if err != nil {
				l.WithError(err).Warnf("Assuming no monsters in map.")
			}

			monstersMax := getMonsterMax(c, len(ableSps))

			toSpawn := monstersMax - monstersInMap
			if toSpawn > 0 {
				result := shuffle(ableSps)
				for i := 0; i < toSpawn; i++ {
					x := result[i]
					processor.CreateMonster(l)(worldId, channelId, mapId, x.Id, x.X, x.Y, x.Fh, x.Team)
				}
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
