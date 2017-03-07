package formatters

type LineCounts struct {
	Missed  int `json:"missed"`
	Covered int `json:"covered"`
	Total   int `json:"total"`
}
