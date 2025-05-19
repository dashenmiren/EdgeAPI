package huaweidns

type CustomLinesResponse struct {
	Lines []struct {
		LineId string `json:"line_id"`
		Name   string `json:"name"`
	} `json:"lines"`
}
