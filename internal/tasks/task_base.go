package tasks

import (
	"github.com/TeaOSLab/EdgeAPI/internal/db/models"
	"github.com/TeaOSLab/EdgeAPI/internal/remotelogs"
)

type BaseTask struct {
}

func (this *BaseTask) logErr(taskType string, errString string) {
	remotelogs.Error("TASK", "run '"+taskType+"' failed: "+errString)
}

func (this *BaseTask) IsPrimaryNode() bool {
	return models.SharedAPINodeDAO.CheckAPINodeIsPrimaryWithoutErr()
}
