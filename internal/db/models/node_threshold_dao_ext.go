//go:build !plus

package models

import "github.com/iwind/TeaGo/dbs"

// FireNodeThreshold 触发相关阈值设置
func (this *NodeThresholdDAO) FireNodeThreshold(tx *dbs.Tx, role string, nodeId int64, item string) error {
	// stub
	return nil
}
