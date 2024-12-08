package bybit

type SMAItem struct {
	Value         float64 `json:"value"`
	Angle         float64 `json:"angle"`
	PreviousAngle float64 `json:"previous_angle"`
}

type SMAData struct {
	FAST  SMAItem `json:"fast"`
	SLOW  SMAItem `json:"slow"`
	HEAVY SMAItem `json:"heavy"`
}

type ActionRecomendation struct {
	Datetime  string  `json:"datetime"`
	Type      string  `json:"type"`
	Candles   int     `json:"candles"`
	Certainty float64 `json:"certainty"`
	Close     float64 `json:"close"`
	Note      string  `json:"note"`
}

type Directions struct {
	Recomendation ActionRecomendation `json:"recomendation"`
	Close         string              `json:"close"`
	Fast          string              `json:"fast_sma"`
	Slow          string              `json:"slow_sma"`
	Heavy         string              `json:"heavy_sma"`
}

type KLineData struct {
	T  int64  `json:"t,omitempty"`
	V  string `json:"v,omitempty"`
	O  string `json:"o,omitempty"`
	C  string `json:"c,omitempty"`
	H  string `json:"h,omitempty"`
	L  string `json:"l,omitempty"`
	S  string `json:"s,omitempty"`
	SN string `json:"sn,omitempty"`

	Datetime       string  `json:"datetime"`
	Timestamp      int64   `json:"timestamp"`
	Symbol         string  `json:"symbol"`
	SymbolInternal string  `json:"symbol_internal"`
	Volume         float64 `json:"volume"`
	Close          float64 `json:"close"`
	Open           float64 `json:"open"`
	High           float64 `json:"high"`
	Low            float64 `json:"low"`

	SMAS       SMAData    `json:"sma"`
	Directions Directions `json:"directions"`
}

type KLineSMAConfig struct {
	Fast  int `json:"fast"`
	Slow  int `json:"slow"`
	Heavy int `json:"heavy"`
}

type KLineParams struct {
	Symbol   string `json:"symbol"`
	Interval string `json:"interval"`
	Limit    int    `json:"limit"`
	To       string `json:"to"`
	From     string `json:"from"`
}

type KLineResolutions struct {
	ResultCount int                   `json:"results"`
	Advices     []ActionRecomendation `json:"advices"`
}

type KLineMeta struct {
	Params KLineParams    `json:"params"`
	SMAS   KLineSMAConfig `json:"smas"`
}

type KLineResponse struct {
	Meta        KLineMeta        `json:"meta"`
	Resolutions KLineResolutions `json:"resolutions"`
	Data        []KLineData      `json:"data"`
}
