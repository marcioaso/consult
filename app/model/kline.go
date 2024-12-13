package model

type KLineAverageConfig struct {
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
	ResultCount int                    `json:"results"`
	Advices     []ActionRecommendation `json:"advices"`
}

type KLineMeta struct {
	Params KLineParams        `json:"params"`
	SMAS   KLineAverageConfig `json:"smas"`
	EMAS   KLineAverageConfig `json:"emas"`
}

type KLineResponse struct {
	Meta        KLineMeta        `json:"meta"`
	Resolutions KLineResolutions `json:"resolutions,omitempty"`
	Data        []KLineData      `json:"data,omitempty"`
}

type KLineData struct {
	Datetime       string  `json:"date_time"`
	Timestamp      int64   `json:"timestamp"`
	Symbol         string  `json:"symbol"`
	SymbolInternal string  `json:"symbol_internal"`
	Volume         float64 `json:"volume"`
	Open           float64 `json:"open"`
	Close          float64 `json:"close"`
	High           float64 `json:"high"`
	Low            float64 `json:"low"`
	Angle          float64 `json:"angle"`

	CloseOpen        float64 `json:"close-open"`
	CloseOpenPercent float64 `json:"close-open_percent"`

	Color string `json:"color"`
}
