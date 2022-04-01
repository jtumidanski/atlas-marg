package reactor

import (
	"atlas-marg/model"
	"atlas-marg/rest/requests"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"strconv"
)

func InMapModelProvider(l logrus.FieldLogger, span opentracing.Span) func(mapId uint32, filters ...model.Filter[Model]) model.SliceProvider[Model] {
	return func(mapId uint32, filters ...model.Filter[Model]) model.SliceProvider[Model] {
		return requests.SliceProvider[attributes, Model](l, span)(requestInMap(mapId), makeModel, filters...)
	}
}

func GetInMap(l logrus.FieldLogger, span opentracing.Span) func(mapId uint32, filters ...model.Filter[Model]) ([]Model, error) {
	return func(mapId uint32, filters ...model.Filter[Model]) ([]Model, error) {
		return InMapModelProvider(l, span)(mapId, filters...)()
	}
}

func makeModel(body requests.DataBody[attributes]) (Model, error) {
	id, err := strconv.ParseUint(body.Id, 10, 32)
	if err != nil {
		return Model{}, err
	}
	attr := body.Attributes
	return Model{
		classification: uint32(id),
		name:           attr.Name,
		x:              attr.X,
		y:              attr.Y,
		delay:          attr.Delay,
		direction:      attr.FacingDirection,
	}, nil
}
