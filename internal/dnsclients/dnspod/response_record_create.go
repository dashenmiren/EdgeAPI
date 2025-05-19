package dnspod

type RecordCreateResponse struct {
	BaseResponse

	Record struct {
		Id     any    `json:"id"`
		Name   string `json:"name"`
		Status string `json:"status"`
	} `json:"record"`
}
