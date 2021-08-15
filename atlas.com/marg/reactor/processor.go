package reactor

import (
	"atlas-marg/map/reactor"
	"github.com/sirupsen/logrus"
	"strconv"
)

func SpawnMissing(l logrus.FieldLogger) func(worldId byte, channelId byte, mapId uint32) {
	return func(worldId byte, channelId byte, mapId uint32) {
		existing, err := GetInMap(l)(worldId, channelId, mapId)
		if err != nil {
			l.WithError(err).Errorf("Unable to retrieve existing reactors in map %d.", mapId)
		}
		needed, err := reactor.GetInMap(l)(mapId)
		for _, nr := range needed {
			if !locationMatch(nr, existing) {
				requestCreate(worldId, channelId, mapId, nr.Classification(), nr.Name(), 0, nr.X(), nr.Y(), nr.Delay(), nr.Direction())
			}
		}

	}
}

func locationMatch(reference reactor.Model, existing []Model) bool {
	for _, er := range existing {
		if er.Classification() == reference.Classification() && er.X() == reference.X() && er.Y() == reference.Y() {
			return true
		}
	}
	return false
}

func GetInMap(l logrus.FieldLogger) func(worldId byte, channelId byte, mapId uint32) ([]Model, error) {
	return func(worldId byte, channelId byte, mapId uint32) ([]Model, error) {
		resp, err := requestInMap(l)(worldId, channelId, mapId)
		if err != nil {
			return nil, err
		}

		reactors := make([]Model, 0)
		for _, d := range resp.Data {
			r, err := makeReactor(d)
			if err != nil {
				l.WithError(err).Errorf("Unable to make reactor %d model.", d.Attributes.Classification)
			} else {
				reactors = append(reactors, *r)
			}
		}
		return reactors, nil
	}
}

func makeReactor(data DataBody) (*Model, error) {
	id, err := strconv.ParseUint(data.Id, 10, 32)
	if err != nil {
		return nil, err
	}
	attr := data.Attributes
	return &Model{
		id:             uint32(id),
		classification: attr.Classification,
		name:           attr.Name,
		state:          attr.State,
		eventState:     attr.EventState,
		delay:          attr.Delay,
		direction:      attr.FacingDirection,
		x:              attr.X,
		y:              attr.Y,
		alive:          attr.Alive,
	}, nil
}
