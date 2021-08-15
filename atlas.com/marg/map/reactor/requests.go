package reactor

import (
	"atlas-marg/rest/requests"
	"fmt"
	"github.com/sirupsen/logrus"
)

const (
	mapInformationServicePrefix string = "/ms/mis/"
	mapInformationService              = requests.BaseRequest + mapInformationServicePrefix
	mapsResource                       = mapInformationService + "maps/"
	reactorsResource                   = mapsResource + "%d/reactors"
)

func requestReactors(l logrus.FieldLogger) func(mapId uint32) (*DataListContainer, error) {
	return func(mapId uint32) (*DataListContainer, error) {
		ar := &DataListContainer{}
		err := requests.Get(l)(fmt.Sprintf(reactorsResource, mapId), ar)
		if err != nil {
			return nil, err
		}
		return ar, nil
	}
}
