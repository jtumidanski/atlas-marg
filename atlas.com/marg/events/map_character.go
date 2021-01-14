package events

type MapCharacterEvent struct {
	WorldId     byte   `json:"worldId"`
	ChannelId   byte   `json:"channelId"`
	MapId       int    `json:"mapId"`
	CharacterId int    `json:"characterId"`
	Type        string `json:"type"`
}
