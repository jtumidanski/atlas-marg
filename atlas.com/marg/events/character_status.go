package events

type CharacterStatusEvent struct {
	WorldId     byte
	ChannelId   byte
	AccountId   int
	CharacterId int
	Type        string
}
