package models

type MapKey struct {
	WorldId   byte
	ChannelId byte
	MapId     int
}

func NewMapKey(worldId byte, channelId byte, mapId int) *MapKey {
	return &MapKey{worldId, channelId, mapId}
}

func (r *MapKey) GetChannelKey() int64 {
	w := int64(int(r.WorldId) * 100000000000)
	c := int64(int(r.ChannelId) * 1000000000)
	return w + c
}

func (r *MapKey) GetMapKey() int64 {
	return r.GetChannelKey() + int64(r.MapId)
}

func GetChannelKey(worldId byte, channelId byte) int64 {
	w := int64(int(worldId) * 100000000000)
	c := int64(int(channelId) * 1000000000)
	return w + c
}

func GetMapKey(worldId byte, channelId byte, mapId int) int64 {
	return GetChannelKey(worldId, channelId) + int64(mapId)
}
