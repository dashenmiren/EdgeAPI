package tasks_test

import (
	"testing"
	"time"

	"github.com/dashenmiren/EdgeAPI/internal/tasks"
	"github.com/iwind/TeaGo/dbs"
)

func TestSSLCertUpdateOCSPTask_Loop(t *testing.T) {
	dbs.NotifyReady()

	var task = tasks.NewSSLCertUpdateOCSPTask(1 * time.Minute)
	err := task.Loop()
	if err != nil {
		t.Fatal(err)
	}
}
