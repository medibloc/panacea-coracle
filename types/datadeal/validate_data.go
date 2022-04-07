package datapool

import (
	"strconv"

	datadealtypes "github.com/medibloc/panacea-core/v2/x/datadeal/types"
)

func NewUnsignedDataValidationCertificate(dealIdStr string, dataHash []byte, encryptedDataUrl []byte, requesterAddress, dataValidatorAddress string) (datadealtypes.UnsignedDataValidationCertificate, error) {
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