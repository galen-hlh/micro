package main

import (
	s "github.com/galen-hlh/micro-sdk/go/idProduce"
	"google.golang.org/grpc/reflection"
	"log"
	"micro/gateway"
	"micro/modules/idProduce"
)

func main() {
	service := gateway.NewService(
		gateway.Name("idProduce"),
	)

	//初始化服务
	service.Init()

	//服务注册
	s.RegisterIdProduceServer(service.Server().GetGrpcServer(), &idProduce.HelperServer{})
	reflection.Register(service.Server().GetGrpcServer())

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
