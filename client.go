package main

import (
	"fmt"
	"github.com/micro/go-micro/config"
	"github.com/micro/go-micro/config/source/consul"
	"github.com/micro/go-micro/util/log"
)

func main() {
	// 注册consul的配置地址
	consulSource := consul.NewSource(
		consul.WithAddress("127.0.0.1:8500"),
		consul.WithPrefix("/micro/config/cluster"),
		// optionally strip the provided prefix from the keys, defaults to false
		consul.StripPrefix(true),
	)
	// 创建新的配置
	conf := config.NewConfig()
	if err := conf.Load(consulSource); err != nil {
		log.Logf("load config errr!!!", err)
	}
	if err := conf.Get("micro", "config", "cluster"); err != nil {
		log.Logf("json format err!!!", err)
	}
	configMap := conf.Map()
	fmt.Println(configMap)
}
