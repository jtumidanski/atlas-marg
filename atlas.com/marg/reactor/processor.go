package reactor

import (
	"atlas-marg/map/reactor"
	"atlas-marg/model"
	"atlas-marg/rest/requests"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"strconv"
)

func SpawnMissing(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32) {
	return func(worldId byte, channelId byte, mapId uint32) {
		existing, err := GetInMap(l, span)(worldId, channelId, mapId)
		if err != nil {
			l.WithError(err).Errorf("Unable to retrieve existing reactors in map %d.", mapId)
		}
		needed, err := reactor.GetInMap(l, span)(mapId, doesNotExist(existing))
		for _, nr := range needed {
			_, _, err = requestCreate(worldId, channelId, mapId, nr.Classification(), nr.Name(), 0, nr.X(), nr.Y(), nr.Delay(), nr.Direction())(l, span)
			if err != nil {
				l.WithError(err).Errorf("Unable to spawn missing reactors.")
			}
		}
	}
}

func doesNotExist(existing []Model) model.Filter[reactor.Model] {
	return func(reference reactor.Model) bool {
		for _, er := range existing {
			if er.Classification() == reference.Classification() && er.X() == reference.X() && er.Y() == reference.Y() {
				return false
			}
		}
		return true
	}
}

func InMapModelProvider(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, filters ...model.Filter[Model]) model.SliceProvider[Model] {
	return func(worldId byte, channelId byte, mapId uint32, filters ...model.Filter[Model]) model.SliceProvider[Model] {
		return requests.SliceProvider[attributes, Model](l, span)(requestInMap(worldId, channelId, mapId), makeModel, filters...)
	}
}

func GetInMap(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, filters ...model.Filter[Model]) ([]Model, error) {
	return func(worldId byte, channelId byte, mapId uint32, filters ...model.Filter[Model]) ([]Model, error) {
		return InMapModelProvider(l, span)(worldId, channelId, mapId, filters...)()
	}
}

func makeModel(data requests.DataBody[attributes]) (Model, error) {
	id, err := strconv.ParseUint(data.Id, 10, 32)
	if err != nil {
		return Model{}, err
	}
	attr := data.Attributes
	return Model{
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
