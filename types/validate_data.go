package types

import (
	panaceatypes "github.com/medibloc/panacea-core/v2/x/market/types"
	"strconv"
)

func NewUnsignedDataValidationCertificate(dealIdStr string, dataHash []byte, encryptedDataUrl []byte, requesterAddress, dataValidatorAddress string) (panaceatypes.UnsignedDataValidationCertificate, error) {
	dealId, err := strconv.ParseUint(dealIdStr, 10, 64)
	if err != nil {
		return panaceatypes.UnsignedDataValidationCertificate{}, err
	}

	return panaceatypes.UnsignedDataValidationCertificate{
		DealId:               dealId,
		DataHash:             dataHash,
		EncryptedDataUrl:     encryptedDataUrl,
		RequesterAddress:     requesterAddress,
		DataValidatorAddress: dataValidatorAddress,
	}, nil
}
