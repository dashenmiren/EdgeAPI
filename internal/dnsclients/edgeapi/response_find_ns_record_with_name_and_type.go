package edgeapi

type FindNSRecordWithNameAndTypeResponse struct {
	BaseResponse

	Data struct {
		NSRecord struct {
			Id       int64  `json:"id"`
			Name     string `json:"name"`
			Type     string `json:"type"`
			Value    string `json:"value"`
			TTL      int32  `json:"ttl"`
			NSRoutes []struct {
				Name string `json:"name"`
				Code string `json:"code"`
			} `json:"nsRoutes"`
		} `json:"nsRecord"`
	}
}
