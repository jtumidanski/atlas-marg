package attributes

type MapCharactersListDataContainer struct {
	Data []MapCharactersData `json:"data"`
}

type MapCharactersData struct {
	Id         string                  `json:"id"`
	Type       string                  `json:"type"`
	Attributes MapCharactersAttributes `json:"attributes"`
}

type MapCharactersAttributes struct {
}
