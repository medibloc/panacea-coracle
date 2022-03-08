package server

import (
	"fmt"

	"github.com/medibloc/panacea-data-market-validator/config"
	"github.com/medibloc/panacea-data-market-validator/types"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Context for data validator application
type Context struct {
	PanaceaConn *grpc.ClientConn
}

func newContext(conf *config.Config) (*Context, error) {
	log.Infof("dial to blockchain: %s", conf.PanaceaGrpcAddress)
	conn, err := grpc.Dial(conf.PanaceaGrpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to blockchain : %w", err)
	}

	return &Context{
		PanaceaConn: conn,
	}, nil
}

func (c Context) Close() error {
	if c.PanaceaConn == nil {
		return types.ErrNoGrpcConnection
	}

	log.Infof("blockchain connection closing")
	return c.PanaceaConn.Close()
}
