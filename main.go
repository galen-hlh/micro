package main

import (
	"google.golang.org/grpc/reflection"
	"log"
	"micro/gateway"
	"micro/gateway/registry"
	h "micro/sdk/go/proto/helper"
	"micro/service/helper"
)

func main() {
	etcdRegistry := registry.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{"127.0.0.1:2379"}
	})
	service := gateway.NewService(
		gateway.Name("helper"),
		gateway.Registry(etcdRegistry),
	)

	//初始化服务
	service.Init()

	//服务注册
	h.RegisterHelperServer(service.Server().GetGrpcServer(), &helper.SnowFlake{})
	reflection.Register(service.Server().GetGrpcServer())

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
