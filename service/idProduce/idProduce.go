package idProduce

import (
	"context"
	"github.com/bwmarrin/snowflake"
	"github.com/galen-hlh/micro-sdk/go/idProduce"
	"micro/util/log"
)

func GetDistributeId(ctx context.Context, request *idProduce.IdProduceRequest) ([]uint64, error) {
	node, err := snowflake.NewNode(1)
	if err != nil {
		return nil, err
	}
	num := request.Len

	var ids []uint64
	if num > 0 {
		for i := uint32(0); i < num; i++ {
			ids = append(ids, uint64(node.Generate()))
		}
	} else {
		ids = append(ids, uint64(node.Generate()))
	}

	log.Logf("Send DistributeId %v", ids)
	return ids, nil
}
