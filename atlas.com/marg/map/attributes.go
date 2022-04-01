package _map

type CharacterDataListContainer struct {
	Data []CharactersDataBody `json:"data"`
}

type CharactersDataBody struct {
	Id         string              `json:"id"`
	Type       string              `json:"type"`
	Attributes CharacterAttributes `json:"attributes"`
}

type CharacterAttributes struct {
}
