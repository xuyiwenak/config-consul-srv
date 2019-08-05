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
		// 这个前缀截止到需要读取的kv前
		consul.WithPrefix("/micro/config/"),
		// 是否使用前缀默认值是关闭，如果不用Get的时候需要填写完整的路径
		consul.StripPrefix(true),
	)
	// 创建新的配置
	conf := config.NewConfig()
	if err := conf.Load(consulSource); err != nil {
		log.Logf("load config errr!!!", err)
	}
	// 可以直接读取到map里面
	configMap := conf.Map()
	fmt.Println(configMap)
	// 也可以根据类型读取
	stmMap := map[string]string{}
	mapConf := conf.Get("cluster").StringMap(stmMap)
	fmt.Printf("%v", mapConf)
}
