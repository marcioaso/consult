package model

type AverageItem struct {
	Value         float64 `json:"value,omitempty"`
	Angle         float64 `json:"angle,omitempty"`
	PreviousAngle float64 `json:"previous_angle,omitempty"`
	Direction     string  `json:"direction,omitempty"`
}

type AverageData struct {
	FAST  AverageItem `json:"fast,omitempty"`
	SLOW  AverageItem `json:"slow,omitempty"`
	HEAVY AverageItem `json:"heavy,omitempty"`
}
