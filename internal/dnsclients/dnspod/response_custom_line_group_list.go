package dnspod

type CustomLineGroupListResponse struct {
	BaseResponse

	Data struct {
		LineGroups []struct {
			Name string `json:"name"`
		} `json:"line_groups"`
		Info struct {
			NowTotal int `json:"now_total"`
			Total    int `json:"total"`
		} `json:"info"`
	} `json:"data"`
}
