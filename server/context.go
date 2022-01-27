package server

import (
	"context"
	"fmt"
	"github.com/medibloc/panacea-data-market-validator/config"
	"github.com/medibloc/panacea-data-market-validator/types"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func newContext(conf *config.Config) (context.Context, error) {
	log.Infof("dial to blockchain: %s", conf.PanaceaGrpcAddress)
	conn, err := grpc.Dial(conf.PanaceaGrpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to blockchain : %w", err)
	}

	return context.WithValue(context.Background(), types.CtxGrpcConnKey, conn), nil
}
