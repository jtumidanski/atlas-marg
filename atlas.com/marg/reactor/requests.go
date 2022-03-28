package reactor

import (
	"atlas-marg/rest/requests"
	"fmt"
)

const (
	reactorServicePrefix string = "/ms/ros/"
	reactorService              = requests.BaseRequest + reactorServicePrefix
	mapReactorsResource         = reactorService + "worlds/%d/channels/%d/maps/%d/reactors"
)

func requestInMap(worldId byte, channelId byte, mapId uint32) requests.Request[attributes] {
	return requests.MakeGetRequest[attributes](fmt.Sprintf(mapReactorsResource, worldId, channelId, mapId))
}

func requestCreate(worldId byte, channelId byte, mapId uint32, classification uint32, name string, state int8, x int16, y int16, delay uint32, direction byte) requests.PostRequest[attributes] {
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
	return requests.MakePostRequest[attributes](fmt.Sprintf(mapReactorsResource, worldId, channelId, mapId), i)
}
