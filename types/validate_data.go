package types

import (
	"strconv"

	datadealtypes "github.com/medibloc/panacea-core/v2/x/datadeal/types"
	datapooltypes "github.com/medibloc/panacea-core/v2/x/datapool/types"
)

func NewUnsignedDataValidationCertificateOfDataDeal(dealIdStr string, dataHash []byte, encryptedDataUrl []byte, requesterAddress, dataValidatorAddress string) (datadealtypes.UnsignedDataValidationCertificate, error) {
	dealId, err := strconv.ParseUint(dealIdStr, 10, 64)
	if err != nil {
		return datadealtypes.UnsignedDataValidationCertificate{}, err
	}

	return datadealtypes.UnsignedDataValidationCertificate{
		DealId:               dealId,
		DataHash:             dataHash,
		EncryptedDataUrl:     encryptedDataUrl,
		RequesterAddress:     requesterAddress,
		DataValidatorAddress: dataValidatorAddress,
	}, nil
}

func NewUnsignedDataValidationCertificateOfDataPool(pool datapooltypes.Pool, dataHash []byte, requesterAddress, dataValidatorAddress string) (datapooltypes.UnsignedDataValidationCertificate, error) {
	return datapooltypes.UnsignedDataValidationCertificate{
		PoolId:        pool.PoolId,
		Round:         pool.Round,
		DataHash:      dataHash,
		DataValidator: dataValidatorAddress,
		Requester:     requesterAddress,
	}, nil
}
