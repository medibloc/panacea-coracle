package server

import (
	"fmt"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	markettypes "github.com/medibloc/panacea-core/v2/x/market/types"

	sdktypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/medibloc/panacea-data-market-validator/config"
	"github.com/medibloc/panacea-data-market-validator/types"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Context for data validator application
type Context struct {
	panaceaConn       *grpc.ClientConn
	interfaceRegistry sdktypes.InterfaceRegistry
}

func newContext(conf *config.Config) (*Context, error) {
	log.Infof("dial to blockchain: %s", conf.PanaceaGrpcAddress)
	conn, err := grpc.Dial(conf.PanaceaGrpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to blockchain : %w", err)
	}

	interfaceRegistry := makeInterfaceRegistry()

	return &Context{
		panaceaConn:       conn,
		interfaceRegistry: interfaceRegistry,
	}, nil
}

func makeInterfaceRegistry() sdktypes.InterfaceRegistry {
	interfaceRegistry := sdktypes.NewInterfaceRegistry()
	authtypes.RegisterInterfaces(interfaceRegistry)
	markettypes.RegisterInterfaces(interfaceRegistry)
	return interfaceRegistry
}

func (c Context) Close() error {
	if c.panaceaConn == nil {
		return types.ErrNoGrpcConnection
	}

	log.Infof("blockchain connection closing")
	return c.panaceaConn.Close()
}
