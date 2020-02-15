package main

import (
	"fmt"
	"micro/grpc"
)

func main() {
	service := grpc.NewService()

	service.Init()

	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
