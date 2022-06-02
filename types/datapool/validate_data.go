package datapool

import (
	datapooltypes "github.com/medibloc/panacea-core/v2/x/datapool/types"
)

func NewUnsignedDataCert(pool datapooltypes.Pool, dataHash []byte, requesterAddress, oracleAddress string) (datapooltypes.UnsignedDataCert, error) {
	return datapooltypes.UnsignedDataCert{
		PoolId:    pool.PoolId,
		Round:     pool.Round,
		DataHash:  dataHash,
		Oracle:    oracleAddress,
		Requester: requesterAddress,
	}, nil
}
