// Copyright 2022 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://cdn.foyeseo.com .

package tasks

type TaskInterface interface {
	Start() error
	Loop() error
	Stop() error
}
