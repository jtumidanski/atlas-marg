package attributes

type MonsterInformationListDataContainer struct {
	Data []MonsterInformationData `json:"data"`
}

type MonsterInformationData struct {
	Id         string                       `json:"id"`
	Type       string                       `json:"type"`
	Attributes MonsterInformationAttributes `json:"attributes"`
}

type MonsterInformationAttributes struct {
	MonsterId int  `json:"monsterId"`
	MobTime   int  `json:"mobTime"`
	Team      int  `json:"team"`
	Cy        int  `json:"cy"`
	F         int  `json:"f"`
	Fh        int  `json:"fh"`
	Rx0       int  `json:"rx0"`
	Rx1       int  `json:"rx1"`
	X         int  `json:"x"`
	Y         int  `json:"y"`
	Hide      bool `json:"hide"`
}
