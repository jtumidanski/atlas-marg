package attributes

type CharacterDataContainer struct {
	Data CharacterData `json:"data"`
}

type CharacterData struct {
	Id         string              `json:"id"`
	Type       string              `json:"type"`
	Attributes CharacterAttributes `json:"attributes"`
}

type CharacterAttributes struct {
	AccountId          int    `json:"accountId"`
	WorldId            int    `json:"worldId"`
	Name               string `json:"name"`
	Level              int    `json:"level"`
	Experience         int    `json:"experience"`
	GachaponExperience int    `json:"gachaponExperience"`
	Strength           int    `json:"strength"`
	Dexterity          int    `json:"dexterity"`
	Luck               int    `json:"luck"`
	Intelligence       int    `json:"intelligence"`
	Hp                 int    `json:"hp"`
	Mp                 int    `json:"mp"`
	MaxHp              int    `json:"maxHp"`
	MaxMp              int    `json:"maxMp"`
	Meso               int    `json:"meso"`
	HpMpUsed           int    `json:"hpMpUsed"`
	JobId              int    `json:"jobId"`
	SkinColor          int    `json:"skinColor"`
	Gender             int    `json:"gender"`
	Fame               int    `json:"fame"`
	Hair               int    `json:"hair"`
	Face               int    `json:"face"`
	Ap                 int    `json:"ap"`
	Sp                 string `json:"sp"`
	MapId              int    `json:"mapId"`
	SpawnPoint         int    `json:"spawnPoint"`
	Gm                 int    `json:"gm"`
	X                  int    `json:"x"`
	Y                  int    `json:"y"`
	Stance             byte   `json:"stance"`
}
