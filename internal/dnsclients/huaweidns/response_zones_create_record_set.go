package huaweidns

type ZonesCreateRecordSetResponse struct {
	Id      string   `json:"id"`
	Line    string   `json:"line"`
	Records []string `json:"records"`
}
