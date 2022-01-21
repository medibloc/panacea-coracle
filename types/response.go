package types

type DataValidationCertificate struct {
	DealId               string
	DataHash             string
	EncryptedDataURL     string
	DataValidatorAddress string
}

type CertificateResponse struct {
	Certificate DataValidationCertificate
	Signature   string
}
