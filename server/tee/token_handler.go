package tee

import (
	"net/http"

	"github.com/medibloc/panacea-data-market-validator/tee"
	log "github.com/sirupsen/logrus"
)

func (svc *teeService) handleToken(writer http.ResponseWriter, request *http.Request) {
	// TODO: Consider creating a Azure attestation token at once when the process is started,
	// rather than whenever HTTP clients call 'GET ../attestation-token'.
	// It's related to the 'exp: 8h' of JWT that Azure sets. It means that we need to handle the recreation of Azure attestation token.
	jwt, err := tee.CreateAzureAttestationToken(svc.TLSCertificate, svc.Conf.EnclaveAttestationProviderURL)
	if err != nil {
		log.Errorf("failed to create Azure attestation token: %v", err)
		http.Error(writer, "failed to create attestation token", http.StatusInternalServerError)
		return
	}
	log.Debugf("Azure attestation token created: %v", jwt)

	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "application/jwt")
	writer.Write([]byte(jwt))
}
