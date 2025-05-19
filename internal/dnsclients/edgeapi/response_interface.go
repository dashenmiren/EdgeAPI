// Copyright 2022 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://cdn.foyeseo.com .

package edgeapi

type ResponseInterface interface {
	IsValid() bool
	Error() error
}
