package helper

import (
	"context"
	"github.com/bwmarrin/snowflake"
	"github.com/galen-hlh/micro-sdk/go/helper"
	"micro/util/log"
)

func GetDistributeId(ctx context.Context, request *helper.IdRequest) (int64, error) {
	node, err := snowflake.NewNode(1)
	if err != nil {
		return 0, err
	}

	id := node.Generate()

	log.Logf("Send distributeId %s", id)
	return int64(id), nil
}
