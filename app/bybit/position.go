package bybit

import (
	"errors"

	"github.com/marcioaso/consult/app/model"
)

type PositionBalance struct {
	Initial       float64                      `json:"total"`
	Realized      float64                      `json:"realized"`
	Close         float64                      `json:"close"`
	Final         float64                      `json:"final"`
	Positions     []model.ActionRecommendation `json:"positions"`
	SoldPositions []model.ActionRecommendation `json:"sold_positions"`
	Multiplier    float64                      `json:"multiplier"`
}

var maxPositions = 2

var Position = PositionBalance{
	Initial: 0,
	Final:   0,
}

func (b *PositionBalance) BuyPosition(price float64, recomendation *model.ActionRecommendation) error {
	multiplier := 1.0
	recomendation.Multiplier = multiplier
	quota := (Position.Final / float64(maxPositions)) * multiplier
	if len(b.Positions) < maxPositions {
		recomendation.Qty = quota / price
		Position.Final -= quota
		Position.Positions = append(Position.Positions, *recomendation)
		return nil
	}
	return errors.New("not enough for quota")
}
func (b *PositionBalance) SellPositions(price float64) (sell []model.ActionRecommendation, profit float64) {
	b.Close = price
	newPositions := make([]model.ActionRecommendation, 0)
	for _, each := range Position.Positions {
		if each.SellAt < price {
			profit = (each.Qty * price) - (each.Qty * each.Close)
			sell = append(Position.SoldPositions, each)
			Position.SoldPositions = append(Position.SoldPositions, each)
			Position.Final += each.Qty * price
		} else {
			newPositions = append(newPositions, each)
		}
	}
	Position.Positions = newPositions
	return
}

func (b *PositionBalance) GetPositions(price float64) (buy []model.ActionRecommendation, sell []model.ActionRecommendation) {
	b.Close = price
	for _, each := range Position.Positions {
		if each.SellAt <= price {
			sell = append(sell, each)
		} else {
			buy = append(buy, each)
		}
	}
	return
}

func (b *PositionBalance) CurrentHold(price float64) (int, float64) {
	hold := 0.0
	boughtPositions := 0
	for _, pos := range Position.Positions {
		if pos.Type == "buy" {
			boughtPositions++
			hold += pos.Qty * price
		}
	}
	return boughtPositions, hold

}
