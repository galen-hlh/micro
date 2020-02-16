package helper

import (
	"context"
	"github.com/bwmarrin/snowflake"
	"micro/sdk/go/proto/helper"
)

type SnowFlake struct{}

func (g *SnowFlake) GetDistributeId(context.Context, *helper.IdRequest) (*helper.IdResponse, error) {
	rsp := new(helper.IdResponse)

	node, err := snowflake.NewNode(1)
	if err != nil {
		return rsp, err
	}

	// Generate a helper ID.
	id := node.Generate()

	rsp.Result = int64(id)
	return rsp, err
}
