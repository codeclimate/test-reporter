package formatters

type LineCounts struct {
	Missed  int `json:"missed"`
	Covered int `json:"covered"`
	Total   int `json:"total"`
}

func (lc LineCounts) CoveredPercent() float64 {
	return (float64(lc.Covered) / float64(lc.Total)) * 100
}
