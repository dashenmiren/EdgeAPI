//go:build !plus

package tasks

// 触发节点动作
func (this *HealthCheckExecutor) fireNodeActions(nodeId int64) error {
	return nil
}
