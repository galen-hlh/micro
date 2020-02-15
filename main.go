package main

import (
	"log"
	"micro/gateway"
)

func main() {
	service := gateway.NewService(
		gateway.Name("helper"),
	)

	service.Init()

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
