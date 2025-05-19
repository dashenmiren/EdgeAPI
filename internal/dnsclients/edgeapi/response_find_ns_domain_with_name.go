package edgeapi

type FindDomainWithNameResponse struct {
	BaseResponse

	Data struct {
		NSDomain struct {
			Id   int64  `json:"id"`
			Name string `json:"name"`
		}
	} `json:"data"`
}
