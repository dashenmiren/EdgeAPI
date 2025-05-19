//go:build !plus

package models

// HasScheduleSettings 检查是否设置了调度
func (this *Node) HasScheduleSettings() bool {
	return false
}
