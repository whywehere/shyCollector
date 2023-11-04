package main

import (
	"flag"
	"fmt"
	"testing"
)

func TestGoFlag(t *testing.T) {
	var (
		host     string
		dbName   string
		port     int
		user     string
		password string
	)
	flag.StringVar(&host, "host", "", "数据库地址")
	flag.StringVar(&dbName, "db_name", "", "数据库名称")
	flag.StringVar(&user, "user", "", "数据库用户")
	flag.StringVar(&password, "password", "", "数据库密码")
	flag.IntVar(&port, "port", 3306, "数据库端口")

	flag.Parse()

	fmt.Printf("数据库地址:%s\n", host)
	fmt.Printf("数据库名称:%s\n", dbName)
	fmt.Printf("数据库用户:%s\n", user)
	fmt.Printf("数据库密码:%s\n", password)
	fmt.Printf("数据库端口:%d\n", port)
}
