package reactor

import (
	"atlas-marg/rest/requests"
	"fmt"
	"github.com/sirupsen/logrus"
)

const (
	reactorServicePrefix string = "/ms/ros/"
	reactorService              = requests.BaseRequest + reactorServicePrefix
	reactorsResource            = reactorService + "reactors"
	mapReactorsResource         = reactorService + "worlds/%d/channels/%d/maps/%d/reactors"
)

func requestInMap(l logrus.FieldLogger) func(worldId byte, channelId byte, mapId uint32) (*DataListContainer, error) {
	return func(worldId byte, channelId byte, mapId uint32) (*DataListContainer, error) {
		dc := &DataListContainer{}
		err := requests.Get(l)(fmt.Sprintf(mapReactorsResource, worldId, channelId, mapId), dc)
		if err != nil {
			return nil, err
		}
		return dc, nil
	}
}

func requestCreate(worldId byte, channelId byte, mapId uint32, classification uint32, name string, state int8, x int16, y int16, delay uint32, direction byte) error {
	i := InputDataContainer{
		Data: DataBody{
			Id:   "",
			Type: "",
			Attributes: Attributes{
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

	_, err := requests.Post(fmt.Sprintf(mapReactorsResource, worldId, channelId, mapId), i)
	if err != nil {
		return err
	}
	return nil
}
