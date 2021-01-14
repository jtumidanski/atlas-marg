package processor

import (
	"atlas-marg/attributes"
	"fmt"
	"log"
	"net/http"
)

type Character struct {
	l *log.Logger
}

func NewCharacter(l *log.Logger) *Character {
	return &Character{l}
}

func (c *Character) GetMapForCharacter(characterId int) (int, error) {
	r, err := http.Get(fmt.Sprintf("http://atlas-nginx:80/ms/cos/characters/%d", characterId))
	if err != nil {
		c.l.Printf("[ERROR] retrieving character data for %d", characterId)
		return 0, err
	}

	td := &attributes.CharacterDataContainer{}
	err = attributes.FromJSON(td, r.Body)
	if err != nil {
		c.l.Printf("[ERROR] decoding character data for %d", characterId)
		return 0, err
	}
	return td.Data.Attributes.MapId, nil
}
