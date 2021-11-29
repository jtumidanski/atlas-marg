package reactor

import (
	"atlas-marg/rest/requests"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

const (
	reactorServicePrefix string = "/ms/ros/"
	reactorService              = requests.BaseRequest + reactorServicePrefix
	reactorsResource            = reactorService + "reactors"
	mapReactorsResource         = reactorService + "worlds/%d/channels/%d/maps/%d/reactors"
)

type Request func(l logrus.FieldLogger, span opentracing.Span) (*dataContainer, error)

func makeRequest(url string) Request {
	return func(l logrus.FieldLogger, span opentracing.Span) (*dataContainer, error) {
		ar := &dataContainer{}
		err := requests.Get(l, span)(url, ar)
		if err != nil {
			return nil, err
		}
		return ar, nil
	}
}

func requestInMap(worldId byte, channelId byte, mapId uint32) Request {
	return makeRequest(fmt.Sprintf(mapReactorsResource, worldId, channelId, mapId))
}

func requestCreate(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, classification uint32, name string, state int8, x int16, y int16, delay uint32, direction byte) error {
	return func(worldId byte, channelId byte, mapId uint32, classification uint32, name string, state int8, x int16, y int16, delay uint32, direction byte) error {
		i := inputDataContainer{
			Data: dataBody{
				Id:   "",
				Type: "",
				Attributes: attributes{
					Classification:  classification,
					Name:            name,
					State:           state,
					X:               x,
					Y:               y,
					Delay:           delay,
					FacingDirection: direction,
				},
			},
		}

		return requests.Post(l, span)(fmt.Sprintf(mapReactorsResource, worldId, channelId, mapId), i, nil, nil)
	}
}
