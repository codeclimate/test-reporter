package formatters

import "math"

type LineCounts struct {
	Missed   int `json:"missed"`
	Covered  int `json:"covered"`
	Total    int `json:"total"`
	Strength int `json:"-"`
}

func (lc LineCounts) CoveredPercent() float64 {
	return (float64(lc.Covered) / float64(lc.Total)) * 100
}

func (lc LineCounts) CoveredStrength() float64 {
	f := float64(lc.Strength) / float64(lc.Total)
	if math.IsNaN(f) {
		return 0
	}
	return f
}
