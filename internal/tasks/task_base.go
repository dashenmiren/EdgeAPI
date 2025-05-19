// Copyright 2022 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://cdn.foyeseo.com .

package tasks

import (
	"github.com/dashenmiren/EdgeAPI/internal/db/models"
	"github.com/dashenmiren/EdgeAPI/internal/remotelogs"
)

type BaseTask struct {
}

func (this *BaseTask) logErr(taskType string, errString string) {
	remotelogs.Error("TASK", "run '"+taskType+"' failed: "+errString)
}

func (this *BaseTask) IsPrimaryNode() bool {
	return models.SharedAPINodeDAO.CheckAPINodeIsPrimaryWithoutErr()
}
