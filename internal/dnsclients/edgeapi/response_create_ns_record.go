// Copyright 2022 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://cdn.foyeseo.com .

package edgeapi

type CreateNSRecordResponse struct {
	BaseResponse

	Data struct {
		NSRecordId int64 `json:"nsRecordId"`
	} `json:"data"`
}
