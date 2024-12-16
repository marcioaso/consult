package model

import (
	"strconv"
	"time"
)

type BybitResponse struct {
	Id int64  `json:"id,omitempty"`
	T  int64  `json:"t,omitempty"`
	V  string `json:"v,omitempty"`
	O  string `json:"o,omitempty"`
	C  string `json:"c,omitempty"`
	H  string `json:"h,omitempty"`
	L  string `json:"l,omitempty"`
	S  string `json:"s,omitempty"`
	SN string `json:"sn,omitempty"`
}

func (kl *BybitResponse) ToData() KLineData {
	kd := KLineData{
		Id:             kl.Id,
		Timestamp:      kl.T,
		Symbol:         kl.S,
		SymbolInternal: kl.SN,
	}
	v, _ := strconv.ParseFloat(kl.V, 64)
	kd.Volume = v

	o, _ := strconv.ParseFloat(kl.O, 64)
	kd.Open = o

	c, _ := strconv.ParseFloat(kl.C, 64)
	kd.Close = c

	h, _ := strconv.ParseFloat(kl.H, 64)
	kd.High = h

	l, _ := strconv.ParseFloat(kl.L, 64)
	kd.Low = l

	dt := time.Unix(kl.T/1000, 0).Format(time.RFC3339)
	kd.Datetime = dt

	kd.CloseOpen = kd.Close - kd.Open
	kd.CloseOpenPercent = (kd.CloseOpen / kd.Open) * 100

	if kd.Close > kd.Open {
		kd.Color = "green"
	} else {
		kd.Color = "red"
	}

	return kd
}
