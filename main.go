package main

import (
	"google.golang.org/grpc/reflection"
	"log"
	"micro/gateway"
	h "micro/sdk/go/proto/helper"
	"micro/service/helper"
)

func main() {
	service := gateway.NewService(
		gateway.Name("helper"),
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
