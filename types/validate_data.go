package types

import "strconv"

func NewUnsignedDataValidationCertificate(dealIdStr string, dataHash []byte, encryptedDataUrl []byte, requesterAddress, dataValidatorAddress string) (UnsignedDataValidationCertificate, error) {
	dealId, err := strconv.ParseUint(dealIdStr, 10, 64)
	if err != nil {
		return UnsignedDataValidationCertificate{}, err
	}

	return UnsignedDataValidationCertificate{
		DealId: dealId,
		DataHash: dataHash,
		EncryptedDataUrl: encryptedDataUrl,
		RequesterAddress: requesterAddress,
		DataValidatorAddress: dataValidatorAddress,
	}, nil
}