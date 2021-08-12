package character

import (
	"github.com/sirupsen/logrus"
)

func GetMapForCharacter(l logrus.FieldLogger) func(characterId uint32) (uint32, error) {
	return func(characterId uint32) (uint32, error) {
		c, err := requestById(l)(characterId)
		if err != nil {
			l.WithError(err).Errorf("Unable to locate the character.")
			return 0, err
		}
		return c.Data.Attributes.MapId, nil
	}
}
