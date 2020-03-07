package main

import (
	h "github.com/galen-hlh/micro-sdk/go/idProduce"
	"google.golang.org/grpc/reflection"
	"log"
	"micro/gateway"
	"micro/modules/idProduce"
)

func main() {
	//etcdRegistry := registry.NewRegistry(func(options *registry.Options) {
	//	options.Addrs = []string{"127.0.0.1:2379"}
	//})
	service := gateway.NewService(
		gateway.Name("helper"),
		gateway.Address("127.0.0.1:9501"),
		//gateway.Registry(etcdRegistry),
	)

	//初始化服务
	service.Init()

	//服务注册
	h.RegisterHelperServer(service.Server().GetGrpcServer(), &idProduce.HelperServer{})
	reflection.Register(service.Server().GetGrpcServer())

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
