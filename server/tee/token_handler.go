package tee

import (
	"net/http"

	"github.com/medibloc/panacea-data-market-validator/tee"
	log "github.com/sirupsen/logrus"
)

func (svc *teeService) handleToken(writer http.ResponseWriter, request *http.Request) {
	// TODO:
	// Consider creating a Azure attestation token at once when the process is started,
	// rather than whenever HTTP clients call 'GET ../attestation-token'.
	// It would be good for reducing the overload of MAA.
	// If so, we must keep in mind that the 'exp' of JWT that MAA sets is 8H.
	// But, the current strategy is not bad in perspective of the MAA overload,
	// since HTTP clients must communicate with MAA to verify JWT anyway.
	jwt, err := tee.CreateAzureAttestationToken(svc.TLSCert.Cert, svc.Conf.EnclaveAttestationProviderURL)
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
