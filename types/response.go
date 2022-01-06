package types

type DataValidationCertificate struct {
	DealId               string
	DataHash             string
	EncryptedDataURL     string
	DataValidatorAddress string
	// TODO sharedECCKey could be added for encryption (ECIES)
}

type CertificateResponse struct {
	Certificate DataValidationCertificate
	Signature   string
}
