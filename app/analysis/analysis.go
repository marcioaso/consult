package analysis

import (
	"github.com/marcioaso/consult/app/model"
	"github.com/marcioaso/consult/pkg"
)

func EnhanceAngleData(data *model.AverageItem, previous model.AverageItem, timeTick float64) *model.AverageItem {
	angle := pkg.GetAngle(
		0,
		previous.Value,
		timeTick,
		data.Value,
	)
	data.Angle = angle
	data.PreviousAngle = previous.Angle

	if angle > 0 {
		data.Direction = "up"
	} else {
		data.Direction = "down"
	}
	return data
}
