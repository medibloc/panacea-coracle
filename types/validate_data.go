package types

import "strconv"

func NewUnsignedDataValidationCertificate(dealIdStr, dataHash, encryptedDataUrl, requesterAddress, dataValidatorAddress string) (UnsignedDataValidationCertificate, error) {
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

func NewDataValidationCertificate(unsignedCert UnsignedDataValidationCertificate, signature []byte) DataValidationCertificate {
	return DataValidationCertificate{
		UnsignedCert: &unsignedCert,
		Signature: signature,
	}
}
