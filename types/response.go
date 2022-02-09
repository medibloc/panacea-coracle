package types

import (
	"encoding/base64"
	panaceatypes "github.com/medibloc/panacea-core/v2/x/market/types"
)

type UnsignedDataValidationCertificateResponse struct {
	DealId                 uint64 `json:"deal_id"`
	DataHashBase64         string `json:"data_hash_base64"`
	EncryptedDataUrlBase64 string `json:"encrypted_data_url_base64"`
	DataValidatorAddress   string `json:"data_validator_address"`
	RequesterAddress       string `json:"requester_address"`
}

type DataValidationCertificateResponse struct {
	UnsignedCert    UnsignedDataValidationCertificateResponse `json:"unsigned_cert"`
	SignatureBase64 string                                    `json:"signature_base64"`
}

// NewUnsignedDataValidationCertificateResponse parse UnsignedDataValidationCertificate
// dataHash and encryptedDataUrl are automatically base64 encoded
func NewUnsignedDataValidationCertificateResponse(certificate panaceatypes.UnsignedDataValidationCertificate) UnsignedDataValidationCertificateResponse {
	return UnsignedDataValidationCertificateResponse{
		DealId:                 certificate.DealId,
		DataHashBase64:         encodeBase64(certificate.DataHash),
		EncryptedDataUrlBase64: encodeBase64(certificate.EncryptedDataUrl),
		DataValidatorAddress:   certificate.DataValidatorAddress,
		RequesterAddress:       certificate.RequesterAddress,
	}
}

// NewDataValidationCertificateResponse parse UnsignedDataValidationCertificate and signature
// signature is automatically base64 encoded
func NewDataValidationCertificateResponse(unsignedCert panaceatypes.UnsignedDataValidationCertificate, signature []byte) DataValidationCertificateResponse {
	return DataValidationCertificateResponse{
		UnsignedCert:    NewUnsignedDataValidationCertificateResponse(unsignedCert),
		SignatureBase64: encodeBase64(signature),
	}
}

func encodeBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}
