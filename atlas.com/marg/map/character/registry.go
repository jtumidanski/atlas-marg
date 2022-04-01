package character

import (
	"errors"
	"sync"
)

type Registry struct {
	mutex sync.Mutex

	mapCharacters map[MapKey][]uint32
	characterMap  map[uint32]MapKey

	mapLocks       map[int64]*sync.Mutex
	characterLocks map[uint32]*sync.Mutex
}

var registry *Registry
var once sync.Once

func GetRegistry() *Registry {
	once.Do(func() {
		registry = &Registry{}
		registry.characterMap = make(map[uint32]MapKey)
		registry.mapCharacters = make(map[MapKey][]uint32)
		registry.mapLocks = make(map[int64]*sync.Mutex)
		registry.characterLocks = make(map[uint32]*sync.Mutex)
	})
	return registry
}

func (r *Registry) getCharacterLock(characterId uint32) *sync.Mutex {
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

func (r *Registry) getMapLock(worldId byte, channelId byte, mapId uint32) *sync.Mutex {
	mk := GetMapKey(worldId, channelId, mapId)
	return r.getMapLockWithKey(mk)
}

func (r *Registry) getMapLockWithKey(mk int64) *sync.Mutex {
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

func remove(c []uint32, i int) []uint32 {
	c[i] = c[len(c)-1]
	return c[:len(c)-1]
}

func indexOf(id uint32, data []uint32) int {
	for k, v := range data {
		if id == v {
			return k
		}
	}
	return -1 //not found.
}

func (r *Registry) removeMapCharacter(mapId MapKey, characterId uint32) {
	index := indexOf(characterId, r.mapCharacters[mapId])
	if index >= 0 && index < len(r.mapCharacters[mapId]) {
		r.mapCharacters[mapId] = remove(r.mapCharacters[mapId], index)
	}
}

func (r *Registry) AddToMap(worldId byte, channelId byte, mapId uint32, characterId uint32) {
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
	mk := NewMapKey(worldId, channelId, mapId)
	if om, ok := r.mapCharacters[*mk]; ok {
		r.mapCharacters[*mk] = append(om, characterId)
	} else {
		r.mapCharacters[*mk] = append([]uint32{}, characterId)
	}
	r.characterMap[characterId] = *mk
	ml.Unlock()
	characterLock.Unlock()
}

func (r *Registry) RemoveFromMap(characterId uint32) {
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

func (r *Registry) GetMapId(characterId uint32) (uint32, error) {
	if mk, ok := r.characterMap[characterId]; ok {
		return mk.MapId, nil
	}
	return 0, errors.New("character not found")
}

func (r *Registry) GetInMap(worldId byte, channelId byte, mapId uint32) []uint32 {
	mk := NewMapKey(worldId, channelId, mapId)
	return r.mapCharacters[*mk]
}

func (r *Registry) GetMapsWithCharacters() []MapKey {
	var result []MapKey
	for i, x := range r.mapCharacters {
		if len(x) > 0 {
			result = append(result, i)
		}
	}
	return result
}
