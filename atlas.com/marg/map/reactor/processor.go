package reactor

import (
	"github.com/sirupsen/logrus"
	"strconv"
)

func GetInMap(l logrus.FieldLogger) func(mapId uint32) ([]Model, error) {
	return func(mapId uint32) ([]Model, error) {
		data, err := requestReactors(l)(mapId)
		if err != nil {
			l.WithError(err).Errorf("Unable to retrieve reactors for map %d.", mapId)
			return make([]Model, 0), err
		}

		results := make([]Model, 0)
		for _, d := range data.Data {
			p, err := makeReactor(d)
			if err != nil {
				return nil, err
			}
			results = append(results, *p)
		}
		return results, nil
	}
}

func makeReactor(body DataBody) (*Model, error) {
	id, err := strconv.ParseUint(body.Id, 10, 32)
	if err != nil {
		return nil, err
	}
	attr := body.Attributes
	return &Model{
		classification: uint32(id),
		name:           attr.Name,
		x:              attr.X,
		y:              attr.Y,
		delay:          attr.Delay,
		direction:      attr.FacingDirection,
	}, nil
}
