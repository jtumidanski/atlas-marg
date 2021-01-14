package events

type MapChangedEvent struct {
	WorldId     byte
	ChannelId   byte
	MapId       int
	PortalId    int
	CharacterId int
}
