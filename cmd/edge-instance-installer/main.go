// Copyright 2024 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://cdn.foyeseo.com .

package main

import (
	"fmt"
	"github.com/dashenmiren/EdgeAPI/internal/instances"
	_ "github.com/iwind/TeaGo/bootstrap"
	"github.com/iwind/TeaGo/lists"
	"log"
	"os"
)

func main() {
	var verbose = lists.ContainsString(os.Args, "-v")

	var dbHost = "127.0.0.1"
	var dbPassword = "123456"
	var dbName = "edges"

	envDBHost, _ := os.LookupEnv("EDGE_DB_HOST")
	if len(envDBHost) > 0 {
		dbHost = envDBHost
		if verbose {
			log.Println("env EDGE_DB_HOST=" + envDBHost)
		}
	}

	envDBPassword, _ := os.LookupEnv("EDGE_DB_PASSWORD")
	if len(envDBPassword) > 0 {
		dbPassword = envDBPassword
		if verbose {
			log.Println("env EDGE_DB_PASSWORD=" + envDBPassword)
		}
	}

	envDBName, _ := os.LookupEnv("EDGE_DB_NAME")
	if len(envDBName) > 0 {
		dbName = envDBName
		if verbose {
			log.Println("env EDGE_DB_NAME=" + envDBName)
		}
	}

	var isTesting = lists.ContainsString(os.Args, "-test") || lists.ContainsString(os.Args, "--test")
	if isTesting {
		fmt.Println("testing mode ...")
	}

	var instance = instances.NewInstance(instances.Options{
		IsTesting: isTesting,
		Verbose:   verbose,
		Cacheable: false,
		WorkDir:   "",
		SrcDir:    "/usr/local/goedge/src",
		DB: struct {
			Host     string
			Port     int
			Username string
			Password string
			Name     string
		}{
			Host:     dbHost,
			Port:     3306,
			Username: "root",
			Password: dbPassword,
			Name:     dbName,
		},
		AdminNode: struct {
			Port int
		}{
			Port: 7788,
		},
		APINode: struct {
			HTTPPort     int
			RestHTTPPort int
		}{
			HTTPPort:     8001,
			RestHTTPPort: 8002,
		},
		Node: struct{ HTTPPort int }{
			HTTPPort: 80,
		},
		UserNode: struct {
			HTTPPort int
		}{
			HTTPPort: 7799,
		},
	})
	err := instance.SetupAll()
	if err != nil {
		fmt.Println("[ERROR]setup failed: " + err.Error())
		return
	}

	fmt.Println("ok")
}
