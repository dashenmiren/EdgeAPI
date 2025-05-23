// Copyright 2022 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://cdn.foyeseo.com .

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
