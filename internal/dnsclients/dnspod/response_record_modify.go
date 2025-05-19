// Copyright 2022 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://cdn.foyeseo.com .

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
