package datapool

//grpcClient This is the Panacea grpcClient used only by pool.
type grpcClient interface {
	GetPubKey(panaceaAddr string) ([]byte, error)
}
