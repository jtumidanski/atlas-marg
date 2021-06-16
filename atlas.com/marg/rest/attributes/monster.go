package attributes

type MonsterInputDataContainer struct {
	Data MonsterData `json:"data"`
}

type MonsterListDataContainer struct {
	Data []MonsterData `json:"data"`
}

type MonsterData struct {
	Id         string            `json:"id"`
	Type       string            `json:"type"`
	Attributes MonsterAttributes `json:"attributes"`
}

type MonsterAttributes struct {
	WorldId            byte            `json:"worldId"`
	ChannelId          byte            `json:"channelId"`
	MapId              int             `json:"mapId"`
	MonsterId          uint32          `json:"monsterId"`
	ControlCharacterId int             `json:"controlCharacterId"`
	X                  int16           `json:"x"`
	Y                  int16           `json:"y"`
	Fh                 uint16          `json:"fh"`
	Stance             int             `json:"stance"`
	Team               int32           `json:"team"`
	Hp                 int             `json:"hp"`
	DamageEntries      []MonsterDamage `json:"damageEntries"`
}

type MonsterDamage struct {
	CharacterId int   `json:"characterId"`
	Damage      int64 `json:"damage"`
}
