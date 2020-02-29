package helper

import (
	"github.com/galen-hlh/micro-sdk/go/helper"
	"go.uber.org/zap"
	HelperService "micro/service/helper"
	"micro/util/log"
)

type HelperServer struct{}

func (g *HelperServer) GetDistributeId(server helper.Helper_GetDistributeIdServer) error {

	request, err := server.Recv()
	if err != nil {
		log.Logf("recv", zap.String("grpc", err.Error()))
		return err
	}

	id, err := HelperService.GetDistributeId(server.Context(), request)
	if err != nil {
		log.Logf("helper", zap.String("get id", err.Error()))
	}

	r := &helper.IdResponse{
		Result: id,
	}

	err = server.Send(r)
	if err != nil {
		log.Logf("send", zap.String("grpc", err.Error()))
		return err
	}

	return nil
	//rsp := new(helper.IdResponse)
	//
	//node, err := snowflake.NewNode(1)
	//if err != nil {
	//	return err
	//}
	//
	//// Generate a helper ID.
	//id := node.Generate()
	//
	//rsp.Result = int64(id)
	//return err
}
