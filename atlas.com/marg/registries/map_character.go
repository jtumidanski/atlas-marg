package registries

import (
	"atlas-marg/models"
	"errors"
	"sync"
)

type MapCharacterRegistry struct {
	mutex sync.Mutex

	mapCharacters map[models.MapKey][]int
	characterMap  map[int]models.MapKey

	mapLocks       map[int64]*sync.Mutex
	characterLocks map[int]*sync.Mutex
}

var mapCharacterRegistry *MapCharacterRegistry
var once sync.Once

func GetMapCharacterRegistry() *MapCharacterRegistry {
	once.Do(func() {
		mapCharacterRegistry = &MapCharacterRegistry{}
		mapCharacterRegistry.characterMap = make(map[int]models.MapKey)
		mapCharacterRegistry.mapCharacters = make(map[models.MapKey][]int)
		mapCharacterRegistry.mapLocks = make(map[int64]*sync.Mutex)
		mapCharacterRegistry.characterLocks = make(map[int]*sync.Mutex)
	})
	return mapCharacterRegistry
}

func (r *MapCharacterRegistry) getCharacterLock(characterId int) *sync.Mutex {
	if val, ok := r.characterLocks[characterId]; ok {
		return val
	} else {
		var cm = &sync.Mutex{}
		r.mutex.Lock()
		r.characterLocks[characterId] = cm
		r.mutex.Unlock()
		return cm
	}
}

func (r *MapCharacterRegistry) getMapLock(worldId byte, channelId byte, mapId int) *sync.Mutex {
	mk := models.GetMapKey(worldId, channelId, mapId)
	return r.getMapLockWithKey(mk)
}

func (r *MapCharacterRegistry) getMapLockWithKey(mk int64) *sync.Mutex {
	if val, ok := r.mapLocks[mk]; ok {
		return val
	} else {
		var mm = &sync.Mutex{}
		r.mutex.Lock()
		r.mapLocks[mk] = mm
		r.mutex.Unlock()
		return mm
	}
}

func remove(c []int, i int) []int {
	c[i] = c[len(c)-1]
	return c[:len(c)-1]
}

func indexOf(id int, data []int) int {
	for k, v := range data {
		if id == v {
			return k
		}
	}
	return -1 //not found.
}

func (r *MapCharacterRegistry) removeMapCharacter(mapId models.MapKey, characterId int) {
	index := indexOf(characterId, r.mapCharacters[mapId])
	if index >= 0 && index < len(r.mapCharacters[mapId]) {
		r.mapCharacters[mapId] = remove(r.mapCharacters[mapId], index)
	}
}

func (r *MapCharacterRegistry) AddCharacterToMap(worldId byte, channelId byte, mapId int, characterId int) {
	characterLock := r.getCharacterLock(characterId)
	characterLock.Lock()
	if om, ok := r.characterMap[characterId]; ok {
		ml := r.getMapLock(worldId, channelId, mapId)
		ml.Lock()
		r.removeMapCharacter(om, characterId)
		ml.Unlock()
	}

	ml := r.getMapLock(worldId, channelId, mapId)
	ml.Lock()
	mk := models.NewMapKey(worldId, channelId, mapId)
	if om, ok := r.mapCharacters[*mk]; ok {
		r.mapCharacters[*mk] = append(om, characterId)
	} else {
		r.mapCharacters[*mk] = append([]int{}, characterId)
	}
	r.characterMap[characterId] = *mk
	ml.Unlock()
	characterLock.Unlock()
}

func (r *MapCharacterRegistry) RemoveCharacterFromMap(characterId int) {
	characterLock := r.getCharacterLock(characterId)
	characterLock.Lock()
	if mk, ok := r.characterMap[characterId]; ok {
		mapLock := r.getMapLockWithKey(mk.GetMapKey())
		mapLock.Lock()
		r.removeMapCharacter(mk, characterId)
		mapLock.Unlock()
		delete(r.characterMap, characterId)
	}
	characterLock.Unlock()
}

func (r *MapCharacterRegistry) GetMapForCharacter(characterId int) (int, error) {
	if mk, ok := r.characterMap[characterId]; ok {
		return mk.MapId, nil
	}
	return 0, errors.New("[WARNING] character not found")
}

func (r *MapCharacterRegistry) GetCharactersInMap(worldId byte, channelId byte, mapId int) []int {
	mk := models.NewMapKey(worldId, channelId, mapId)
	return r.mapCharacters[*mk]
}

func (r *MapCharacterRegistry) GetMapsWithCharacters() []models.MapKey {
	var result []models.MapKey
	for i, x := range r.mapCharacters {
		if len(x) > 0 {
			result = append(result, i)
		}
	}
	return result
}
