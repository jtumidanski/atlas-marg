package character

import (
	"atlas-marg/rest/requests"
	"fmt"
	"github.com/sirupsen/logrus"
)

const (
	charactersServicePrefix string = "/ms/cos/"
	charactersService              = requests.BaseRequest + charactersServicePrefix
	charactersResource             = charactersService + "characters"
	charactersById                 = charactersResource + "/%d"
)

func requestById(l logrus.FieldLogger) func(characterId uint32) (*DataContainer, error) {
	return func(characterId uint32) (*DataContainer, error) {
		ar := &DataContainer{}
		err := requests.Get(l)(fmt.Sprintf(charactersById, characterId), ar)
		if err != nil {
			return nil, err
		}
		return ar, nil
	}
}