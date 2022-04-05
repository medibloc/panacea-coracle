package datapool

import (
	datapooltypes "github.com/medibloc/panacea-core/v2/x/datapool/types"
)

func NewUnsignedDataValidationCertificate(pool datapooltypes.Pool, dataHash []byte, requesterAddress, dataValidatorAddress string) (datapooltypes.UnsignedDataValidationCertificate, error) {
	return datapooltypes.UnsignedDataValidationCertificate{
		PoolId:        pool.PoolId,
		Round:         pool.Round,
		DataHash:      dataHash,
		DataValidator: dataValidatorAddress,
		Requester:     requesterAddress,
	}, nil
}
