package xccov

type sourceFile struct {
	Path  string `json:"path"`
	Functions []struct {
		CoveredLines   int `json:"coveredLines"`
		LineNumber int `json:"lineNumber"`
		ExecutableLines int `json:"executableLines"`
	} `json:"functions"`
}

type xccovFile struct {
	Targets []struct {
		Files []sourceFile `json:"files"`
	} `json:"targets"`
}
