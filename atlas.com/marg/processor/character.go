package processor

import (
	"atlas-marg/rest/attributes"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
)

func GetMapForCharacter(l logrus.FieldLogger) func(characterId int) (int, error) {
	return func(characterId int) (int, error) {
		r, err := http.Get(fmt.Sprintf("http://atlas-nginx:80/ms/cos/characters/%d", characterId))
		if err != nil {
			l.WithError(err).Errorf("Retrieving character data for %d", characterId)
			return 0, err
		}

		td := &attributes.CharacterDataContainer{}
		err = attributes.FromJSON(td, r.Body)
		if err != nil {
			l.WithError(err).Errorf("Decoding character data for %d", characterId)
			return 0, err
		}
		return td.Data.Attributes.MapId, nil
	}
}
