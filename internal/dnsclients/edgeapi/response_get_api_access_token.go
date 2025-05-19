package edgeapi

type GetAPIAccessToken struct {
	BaseResponse

	Data struct {
		Token     string `json:"token"`
		ExpiresAt int64  `json:"expiresAt"`
	} `json:"data"`
}
