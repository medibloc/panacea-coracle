package datapool

import (
	"fmt"
	"github.com/medibloc/panacea-data-market-validator/config"
	"github.com/medibloc/panacea-data-market-validator/crypto"
	"github.com/medibloc/panacea-data-market-validator/panacea"
	"github.com/medibloc/panacea-data-market-validator/server/service"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func makeAccount() (*panacea.ValidatorAccount, error) {
	mnemonic, err := crypto.GenerateMnemonic()

	acc, err := panacea.NewValidatorAccount(mnemonic)
	if err != nil {
		return nil, err
	}
	return acc, nil
}

func TestHandleDownloadData(t *testing.T) {
	requesterAcc, err := makeAccount()
	require.NoError(t, err)

	// Make an HTTP request and an HTTP server simulator (recorder)
	req := httptest.NewRequest(http.MethodGet, "/v0/tee/attestation-token", nil)
	req.Header.Add("Content-Type", "application/octet-stream")
	req.Form = make(map[string][]string)
	req.Form.Add("requester_address", requesterAcc.GetAddress())
	recorder := httptest.NewRecorder()

	homeDir, err := os.UserHomeDir()
	require.NoError(t, err)
	confDir := filepath.Join(homeDir, ".dataval")
	conf, err := config.ReadConfigTOML(filepath.Join(confDir, "config.toml"))

	baseSvc, err := service.New(conf)
	require.NoError(t, err)
	svc := dataPoolService{
		baseSvc,
	}

	svc.handleDownloadData(recorder, req)

}

func TestHandleDownloadData2(t *testing.T) {
	ch := make(chan int)

	fmt.Println(0)
	go func() {
		ch <- 1
		fmt.Println("input 1")
		ch <- 2
		fmt.Println("input 2")
		ch <- 3
		fmt.Println("input 3")
		ch <- 4
		fmt.Println("input 4")
		ch <- 5
		fmt.Println("input 5")
	}()

	fmt.Println(<-ch)
	time.Sleep(time.Second)
	fmt.Println(<-ch)
	time.Sleep(time.Second)
	fmt.Println(<-ch)
	time.Sleep(time.Second)
	fmt.Println(<-ch)
}
