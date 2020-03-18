# micro
golang简单实现grpc微服务，部分代码借鉴于go-micro

## 项目包管理基于go mod管理需要配置
```bash
GO111MODULE=on
```

## demo
```golang
package main

import (
	s "github.com/galen-hlh/micro-sdk/go/idProduce"
	"google.golang.org/grpc/reflection"
	"log"
	"micro/gateway"
	"micro/gateway/registry"
	"micro/modules/idProduce"
)

func main() {
    //设置etcd集群
	etcdRegistry := registry.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{"127.0.0.1:2379"}
	})
	service := gateway.NewService(
		gateway.Name("idProduce"),            //设置微服务名称
		gateway.Address("127.0.0.1:9501"),    //设置对外暴露端口
		gateway.Registry(etcdRegistry),       //设置服务使用etcd注册
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

```

## 运行项目
```bash
go run main.go
```
