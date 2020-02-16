package helper

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"micro/gateway/server"
	"micro/sdk/go/proto/helper"
	"testing"
)

func TestSnowFlake_GetDistributeId(t *testing.T) {
	// 建立连接到gRPC服务
	conn, err := grpc.Dial(server.DefaultAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	// 函数结束时关闭连接
	defer conn.Close()

	// 创建Waiter服务的客户端
	client := helper.NewHelperClient(conn)

	// 调用gRPC接口
	tr, err := client.GetDistributeId(context.Background(), &helper.IdRequest{Code: 1})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("服务端响应: %d", tr.Result)
}
