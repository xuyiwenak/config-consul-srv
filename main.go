package main

import (
	"github.com/micro/go-config"
	"github.com/micro/go-config/source/consul"
	"github.com/micro/go-log"
)

type dbConfig struct {
	DB int `json:"dict"`
	//Ca int    `json:"port"`
}

func main() {
	// 注册consul的配置地址
	consulSource := consul.NewSource(
		consul.WithAddress("127.0.0.1:8500"),
		consul.WithPrefix("/micro/config/database"),
		// optionally strip the provided prefix from the keys, defaults to false
		consul.StripPrefix(true),
	)
	// 创建新的配置
	conf := config.NewConfig()
	if err := conf.Load(consulSource); err != nil {
		log.Logf("load config errr!!!", err)
	}
	var d dbConfig
	aa := conf.Map()
	log.Log(aa)
	if err := conf.Get("micro", "config", "database", "loop").Scan(&d); err != nil {
		log.Logf("json format err!!!", err)
	}
	log.Log(d.DB)
}
