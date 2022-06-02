package datapool

import (
	"strconv"

	datadealtypes "github.com/medibloc/panacea-core/v2/x/datadeal/types"
)

func NewUnsignedDataCert(dealIdStr string, dataHash []byte, encryptedDataUrl []byte, requesterAddress, oracleAddress string) (datadealtypes.UnsignedDataCert, error) {
	dealId, err := strconv.ParseUint(dealIdStr, 10, 64)
	if err != nil {
		return datadealtypes.UnsignedDataCert{}, err
	}

	return datadealtypes.UnsignedDataCert{
		DealId:           dealId,
		DataHash:         dataHash,
		EncryptedDataUrl: encryptedDataUrl,
		RequesterAddress: requesterAddress,
		OracleAddress:    oracleAddress,
	}, nil
}
