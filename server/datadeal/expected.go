package datadeal

import markettypes "github.com/medibloc/panacea-core/v2/x/market/types"

//grpcClient This is the Panacea grpcClient used only by deal.
type grpcClient interface {
	GetPubKey(panaceaAddr string) ([]byte, error)
	GetDeal(id string) (markettypes.Deal, error)
}
