package dnspod

type RecordModifyResponse struct {
	BaseResponse

	Record struct {
		Id     any    `json:"id"`
		Name   string `json:"name"`
		Value  string `json:"value"`
		Status string `json:"status"`
	} `json:"record"`
}
