package _map

import (
	"atlas-marg/map/character"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"sync"
)

var dungeons []DungeonConfig
var once sync.Once

func Dungeons() []DungeonConfig {
	once.Do(func() {
		dungeons = make([]DungeonConfig, 0)
		dungeons = append(dungeons, DungeonConfig{dungeonId: 105050101, amount: 30}) // cave of mushrooms
		dungeons = append(dungeons, DungeonConfig{dungeonId: 105040320, amount: 34}) // golem castle ruins
		dungeons = append(dungeons, DungeonConfig{dungeonId: 260020630, amount: 30}) // hill of sandstorms
		dungeons = append(dungeons, DungeonConfig{dungeonId: 100020100, amount: 30}) // henesys pig farm
		dungeons = append(dungeons, DungeonConfig{dungeonId: 105090320, amount: 30}) // drakes blue cave
		dungeons = append(dungeons, DungeonConfig{dungeonId: 221023401, amount: 30}) // drummer bunnys lair
		dungeons = append(dungeons, DungeonConfig{dungeonId: 240020512, amount: 30}) // the round table of kentarus
		dungeons = append(dungeons, DungeonConfig{dungeonId: 240040800, amount: 19}) // the restoring memory
		dungeons = append(dungeons, DungeonConfig{dungeonId: 240040900, amount: 19}) // newt secured zone
		dungeons = append(dungeons, DungeonConfig{dungeonId: 251010410, amount: 30}) // pillage of treasure island
		dungeons = append(dungeons, DungeonConfig{dungeonId: 261020301, amount: 30}) // critical error
		dungeons = append(dungeons, DungeonConfig{dungeonId: 551030001, amount: 19}) // longest ride on byebye station
	})
	return dungeons
}

func InMiniDungeon(l logrus.FieldLogger, _ opentracing.Span) func(characterId uint32) (bool, error) {
	return func(characterId uint32) (bool, error) {
		mapId, err := character.GetRegistry().GetMapId(characterId)
		if err != nil {
			l.WithError(err).Errorf("Unable to locate character %d.", characterId)
			return false, err
		}
		return IsMiniDungeon(mapId), nil
	}
}

func IsMiniDungeon(mapId uint32) bool {
	for _, dc := range Dungeons() {
		if mapId >= dc.dungeonId && mapId <= dc.dungeonId+uint32(dc.amount) {
			return true
		}
	}
	return false
}

type DungeonConfig struct {
	dungeonId uint32
	amount    uint8
}
