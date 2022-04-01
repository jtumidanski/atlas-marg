package monster

import (
	"atlas-marg/map/character"
	"atlas-marg/monster"
	"atlas-marg/rest/requests"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"math"
	"math/rand"
	"time"
)

func Spawn(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32) {
	return func(worldId byte, channelId byte, mapId uint32) {
		c := len(character.GetRegistry().GetInMap(worldId, channelId, mapId))
		if c > 0 {
			sps, err := GetSpawnPoints(l, span)(mapId)
			if err != nil {
				l.WithError(err).Errorf("Unable to get spawn points for map %d.", mapId)
				return
			}

			var ableSps []SpawnPoint
			for _, x := range sps {
				if x.MobTime >= 0 {
					ableSps = append(ableSps, x)
				}
			}

			monstersInMap, err := monster.CountInMap(l, span)(worldId, channelId, mapId)
			if err != nil {
				l.WithError(err).Warnf("Assuming no monsters in map.")
			}

			monstersMax := getMonsterMax(c, len(ableSps))

			toSpawn := monstersMax - monstersInMap
			if toSpawn > 0 {
				result := shuffle(ableSps)
				for i := 0; i < toSpawn; i++ {
					x := result[i]
					monster.CreateMonster(l, span)(worldId, channelId, mapId, x.Id, x.X, x.Y, x.Fh, x.Team)
				}
			}
		}
	}
}

func shuffle(vals []SpawnPoint) []SpawnPoint {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	ret := make([]SpawnPoint, len(vals))
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

func GetSpawnPoints(l logrus.FieldLogger, span opentracing.Span) func(mapId uint32) ([]SpawnPoint, error) {
	return func(mapId uint32) ([]SpawnPoint, error) {
		data, err := requestSpawnPoints(mapId)(l, span)
		if err != nil {
			return nil, err
		}
		var result []SpawnPoint
		for _, d := range data.DataList() {
			result = append(result, makeModel(d))
		}
		return result, nil
	}
}

func makeModel(data requests.DataBody[attributes]) SpawnPoint {
	attr := data.Attributes
	return SpawnPoint{
		Id:      attr.Id,
		MobTime: attr.MobTime,
		Team:    attr.Team,
		Cy:      attr.CY,
		F:       attr.F,
		Fh:      attr.FH,
		Rx0:     attr.RX0,
		Rx1:     attr.RX1,
		X:       attr.X,
		Y:       attr.Y,
	}
}
