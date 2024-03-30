package edgeapi

type FindAllNSRoutesResponse struct {
	BaseResponse

	Data struct {
		NSRoutes []struct {
			Name string `json:"name"`
			Code string `json:"code"`
		} `json:"nsRoutes"`
	} `json:"data"`
}
