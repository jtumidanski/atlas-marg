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
	WorldId            int            `json:"worldId"`
	ChannelId          int            `json:"channelId"`
	MapId              int             `json:"mapId"`
	MonsterId          int             `json:"monsterId"`
	ControlCharacterId int             `json:"controlCharacterId"`
	X                  int             `json:"x"`
	Y                  int             `json:"y"`
	Fh                 int             `json:"fh"`
	Stance             int             `json:"stance"`
	Team               int             `json:"team"`
	Hp                 int             `json:"hp"`
	DamageEntries      []MonsterDamage `json:"damageEntries"`
}

type MonsterDamage struct {
	CharacterId int   `json:"characterId"`
	Damage      int64 `json:"damage"`
}
