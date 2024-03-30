package cloudflare

type ZonesResponse struct {
	BaseResponse

	Result []struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"result"`
}
