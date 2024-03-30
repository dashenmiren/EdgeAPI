package dnspod

type RecordLineResponse struct {
	BaseResponse

	Lines []string `json:"lines"`
}
