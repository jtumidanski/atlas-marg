package character

import (
	"atlas-marg/rest/requests"
	"fmt"
)

const (
	charactersServicePrefix string = "/ms/cos/"
	charactersService              = requests.BaseRequest + charactersServicePrefix
	charactersResource             = charactersService + "characters"
	charactersById                 = charactersResource + "/%d"
)

func requestById(characterId uint32) requests.Request[Attributes] {
	return requests.MakeGetRequest[Attributes](fmt.Sprintf(charactersById, characterId))
}
