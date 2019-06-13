package main

import (
	"github.com/micro/go-log"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/consul"
	"time"
)

func main() {
	Init()
	// 注册一个consul服务
	micReg := consul.NewRegistry(registryOptions)

	// 新建服务
	service := micro.NewService(
		micro.Name("bambooRat.micro.srv.config"),
		micro.Registry(micReg),
		micro.Version("latest"),
	)

	// 服务初始化
	service.Init()

	// 启动服务
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

func registryOptions(ops *registry.Options) {
	ops.Timeout = time.Second * 5
	ops.Addrs = []string{consulConfigCenterAddr}
}
