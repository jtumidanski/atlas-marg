package reactor

import (
	"atlas-marg/rest/requests"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"strconv"
)

type ModelOperator func(*Model)

type ModelListOperator func([]*Model)

type ModelProvider func() (*Model, error)

type ModelListProvider func() ([]*Model, error)

type Filter func(*Model) bool

func requestModelListProvider(l logrus.FieldLogger, span opentracing.Span) func(r requests.Request[attributes], filters ...Filter) ModelListProvider {
	return func(r requests.Request[attributes], filters ...Filter) ModelListProvider {
		return func() ([]*Model, error) {
			resp, err := r(l, span)
			if err != nil {
				return nil, err
			}

			ms := make([]*Model, 0)
			for _, v := range resp.DataList() {
				m, err := makeModel(v)
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

func InMapModelProvider(l logrus.FieldLogger, span opentracing.Span) func(mapId uint32, filters ...Filter) ModelListProvider {
	return func(mapId uint32, filters ...Filter) ModelListProvider {
		return requestModelListProvider(l, span)(requestInMap(mapId), filters...)
	}
}

func GetInMap(l logrus.FieldLogger, span opentracing.Span) func(mapId uint32, filters ...Filter) ([]*Model, error) {
	return func(mapId uint32, filters ...Filter) ([]*Model, error) {
		return InMapModelProvider(l, span)(mapId, filters...)()
	}
}

func makeModel(body requests.DataBody[attributes]) (*Model, error) {
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
