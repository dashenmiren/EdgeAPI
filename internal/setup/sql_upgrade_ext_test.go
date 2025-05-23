// Copyright 2022 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://cdn.foyeseo.com .
//go:build !plus

package setup

import (
	"github.com/iwind/TeaGo/dbs"
	"testing"
)

func TestUpgradeSQLData_v0_5_6(t *testing.T) {
	db, err := dbs.NewInstanceFromConfig(&dbs.DBConfig{
		Driver: "mysql",
		Dsn:    "root:123456@tcp(127.0.0.1:3306)/db_edge?charset=utf8mb4&timeout=30s",
		Prefix: "edge",
	})
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = db.Close()
	}()

	err = upgradeV0_5_6(db)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("ok")
}


func TestUpgradeSQLData_v1_3_4(t *testing.T) {
	db, err := dbs.NewInstanceFromConfig(&dbs.DBConfig{
		Driver: "mysql",
		Dsn:    "root:123456@tcp(127.0.0.1:3306)/db_edge?charset=utf8mb4&timeout=30s",
		Prefix: "edge",
	})
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = db.Close()
	}()

	err = upgradeV1_3_4(db)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("ok")
}


