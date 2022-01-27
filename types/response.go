package types

import "encoding/base64"

type UnsignedDataValidationCertificateResponse struct {
	DealId                 uint64
	DataHashBase64         string
	EncryptedDataUrlBase64 string
	DataValidatorAddress   string
	RequesterAddress       string
}

type DataValidationCertificateResponse struct {
	UnsignedCert    UnsignedDataValidationCertificateResponse
	SignatureBase64 string
}

// NewUnsignedDataValidationCertificateResponse parse UnsignedDataValidationCertificate
// dataHash and encryptedDataUrl are automatically base64 encoded
func NewUnsignedDataValidationCertificateResponse(certificate UnsignedDataValidationCertificate) UnsignedDataValidationCertificateResponse {
	return UnsignedDataValidationCertificateResponse{
		DealId: certificate.DealId,
		DataHashBase64:         encodeBase64(certificate.DataHash),
		EncryptedDataUrlBase64: encodeBase64(certificate.EncryptedDataUrl),
		DataValidatorAddress:   certificate.DataValidatorAddress,
		RequesterAddress:       certificate.RequesterAddress,
	}
}

// NewDataValidationCertificateResponse parse UnsignedDataValidationCertificate and signature
// signature is automatically base64 encoded
func NewDataValidationCertificateResponse(unsignedCert UnsignedDataValidationCertificate, signature []byte) DataValidationCertificateResponse {
	return DataValidationCertificateResponse{
		UnsignedCert: NewUnsignedDataValidationCertificateResponse(unsignedCert),
		SignatureBase64: encodeBase64(signature),
	}
}

func encodeBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

