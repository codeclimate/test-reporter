package formatters

import "math"

type LineCounts struct {
	Missed  int `json:"missed"`
	Covered int `json:"covered"`
	Total   int `json:"total"`
}

func (lc LineCounts) CoveredPercent() float64 {
	f := (float64(lc.Covered) / float64(lc.Total)) * 100
	if math.IsNaN(f) {
		return 0
	}
	return f
}
