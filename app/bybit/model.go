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

type Directions struct {
	Close string `json:"close"`
	Fast  string `json:"fast_sma"`
	Slow  string `json:"slow_sma"`
	Heavy string `json:"heavy_sma"`
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

type KLineDefinitions struct {
	SMAS KLineSMAConfig `json:"smas"`
}

type KLineResponse struct {
	Definitions KLineDefinitions `json:"definitions"`
	Data        []KLineData      `json:"data"`
}
