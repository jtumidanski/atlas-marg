package reactor

import (
	"atlas-marg/map/reactor"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"strconv"
)

type ModelOperator func(*Model)

type ModelListOperator func([]*Model)

type ModelProvider func() (*Model, error)

type ModelListProvider func() ([]*Model, error)

type Filter func(*Model) bool

func SpawnMissing(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32) {
	return func(worldId byte, channelId byte, mapId uint32) {
		existing, err := GetInMap(l, span)(worldId, channelId, mapId)
		if err != nil {
			l.WithError(err).Errorf("Unable to retrieve existing reactors in map %d.", mapId)
		}
		needed, err := reactor.GetInMap(l, span)(mapId, doesNotExist(existing))
		for _, nr := range needed {
			err = requestCreate(l, span)(worldId, channelId, mapId, nr.Classification(), nr.Name(), 0, nr.X(), nr.Y(), nr.Delay(), nr.Direction())
			if err != nil {
				l.WithError(err).Errorf("Unable to spawn missing reactors.")
			}
		}
	}
}

func doesNotExist(existing []*Model) reactor.Filter {
	return func(reference *reactor.Model) bool {
		for _, er := range existing {
			if er.Classification() == reference.Classification() && er.X() == reference.X() && er.Y() == reference.Y() {
				return false
			}
		}
		return true
	}
}

func requestModelListProvider(l logrus.FieldLogger, span opentracing.Span) func(r Request, filters ...Filter) ModelListProvider {
	return func(r Request, filters ...Filter) ModelListProvider {
		return func() ([]*Model, error) {
			resp, err := r(l, span)
			if err != nil {
				return nil, err
			}

			ms := make([]*Model, 0)
			for _, v := range resp.DataList() {
				m, err := makeModel(&v)
				if err != nil {
					return nil, err
				}
				ok := true
				for _, filter := range filters {
					if !filter(m) {
						ok = false
						break
					}
				}
				if ok {
					ms = append(ms, m)
				}
			}
			return ms, nil
		}
	}
}

func InMapModelProvider(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, filters ...Filter) ModelListProvider {
	return func(worldId byte, channelId byte, mapId uint32, filters ...Filter) ModelListProvider {
		return requestModelListProvider(l, span)(requestInMap(worldId, channelId, mapId), filters...)
	}
}

func GetInMap(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, filters ...Filter) ([]*Model, error) {
	return func(worldId byte, channelId byte, mapId uint32, filters ...Filter) ([]*Model, error) {
		return InMapModelProvider(l, span)(worldId, channelId, mapId, filters...)()
	}
}

func makeModel(data *dataBody) (*Model, error) {
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
