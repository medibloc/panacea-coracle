package types

type DataValidationCertificate struct {
	DealId               string
	DataHash             string
	EncryptedDataURL     string
	DataValidatorAddress string
	// TODO sharedECCKey (for ECIES)가 포함될 수 있음
}

type CertificateResponse struct {
	Certificate DataValidationCertificate
	Signature   string
}
