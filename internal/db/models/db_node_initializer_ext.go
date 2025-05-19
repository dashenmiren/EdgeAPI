//go:build !plus

package models

import "github.com/iwind/TeaGo/dbs"

var nsAccessLogDAOMapping = map[int64]any{} // dbNodeId => DAO

func initAccessLogDAO(db *dbs.DB, node *DBNode) {
}
