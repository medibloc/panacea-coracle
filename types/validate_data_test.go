package types_test

import (
	"fmt"
	"github.com/medibloc/panacea-data-market-validator/types"
	"testing"
)

func TestVal(t *testing.T) {
	cer := types.UnsignedDataValidationCertificate{
		DealId: 1,
		EncryptedDataUrl: "asdf41",

	}

	fmt.Println(cer.Marshal())
}