package dnspod

type DomainInfoResponse struct {
	BaseResponse

	Domain struct {
		Id    any    `json:"id"`
		Name  string `json:"name"`
		Grade string `json:"grade"`
	} `json:"domain"`
}
