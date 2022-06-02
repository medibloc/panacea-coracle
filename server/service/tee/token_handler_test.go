package tee

import (
	"encoding/binary"
	"encoding/hex"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/edgelesssys/ego/attestation"
	"github.com/medibloc/panacea-oracle/config"
	"github.com/medibloc/panacea-oracle/server/service"
	"github.com/medibloc/panacea-oracle/tee"
	"github.com/stretchr/testify/require"
)

func TestHandleToken(t *testing.T) {
	// Load env vars for testing
	enclaveSignerID, err := hex.DecodeString(os.Getenv("EDG_TEST_ENCLAVE_SIGNER_ID_HEX"))
	require.NoError(t, err)

	// Make an HTTP request and an HTTP server simulator (recorder)
	req := httptest.NewRequest(http.MethodGet, "/v0/tee/attestation-token", nil)
	recorder := httptest.NewRecorder()

	// Prepare a service struct and execute the HTTP request
	tlsCert, err := tee.CreateTLSCertificate()
	require.NoError(t, err)
	enclaveConfig := config.EnclaveConfig{Enable: true, AttestationProviderAddr: "https://shareduks.uks.attest.azure.net"}
	svc := &teeService{
		&service.Service{
			Conf:    &config.Config{Enclave: enclaveConfig},
			TLSCert: tlsCert,
		},
	}
	svc.handleToken(recorder, req)

	// Get the HTTP response and verify it
	res := recorder.Result()
	defer res.Body.Close()

	require.Equal(t, http.StatusOK, res.StatusCode)

	attestationTokenBytes, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)

	report, err := attestation.VerifyAzureAttestationToken(string(attestationTokenBytes), svc.Conf.Enclave.AttestationProviderAddr)
	require.NoError(t, err)
	t.Log("Azure attestation token verified")

	// Verify report values with that were defined in the enclave.json
	// and that were included into the test binary during build.
	require.Equal(t, []byte(enclaveSignerID), report.SignerID)
	require.Equal(t, uint16(1), binary.LittleEndian.Uint16(report.ProductID))
	require.Equal(t, uint(1), report.SecurityVersion)
	t.Log("Attestation report verified")
}
