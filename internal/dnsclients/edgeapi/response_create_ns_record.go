package edgeapi

type CreateNSRecordResponse struct {
	BaseResponse

	Data struct {
		NSRecordId int64 `json:"nsRecordId"`
	} `json:"data"`
}
