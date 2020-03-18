package idProduce

import (
	"github.com/galen-hlh/micro-sdk/go/idProduce"
	"go.uber.org/zap"
	service "micro/service/idProduce"
	"micro/util/log"
)

type HelperServer struct{}

func (g *HelperServer) GetDistributeId(server idProduce.IdProduce_GetDistributeIdServer) error {

	request, err := server.Recv()
	if err != nil {
		log.Logf("recv", zap.String("grpc", err.Error()))
		return err
	}

	ids, err := service.GetDistributeId(server.Context(), request)
	if err != nil {
		log.Logf("helper", zap.String("get id", err.Error()))
	}

	r := &idProduce.IdProduceResponse{
		Ids: ids,
	}

	err = server.Send(r)
	if err != nil {
		log.Logf("send", zap.String("grpc", err.Error()))
		return err
	}

	return nil
}
