package accounts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ElrondNetwork/elastic-indexer-go/converters"
	"github.com/ElrondNetwork/elastic-indexer-go/data"
)

const (
	esdtIdentifierSeparator  = "-"
	esdtRandomSequenceLength = 6
)

// SerializeNFTCreateInfo will serialize the provided nft create information in a way that Elastic Search expects a bulk request
func (ap *accountsProcessor) SerializeNFTCreateInfo(tokensInfo []*data.TokenInfo) ([]*bytes.Buffer, error) {
	buffSlice := data.NewBufferSlice()
	for _, tokenData := range tokensInfo {
		meta := []byte(fmt.Sprintf(`{ "index" : { "_id" : "%s" } }%s`, tokenData.Identifier, "\n"))
		serializedData, errMarshal := json.Marshal(tokenData)
		if errMarshal != nil {
			return nil, errMarshal
		}

		err := buffSlice.PutData(meta, serializedData)
		if err != nil {
			return nil, err
		}
	}

	return buffSlice.Buffers(), nil
}

// SerializeAccounts will serialize the provided accounts in a way that Elastic Search expects a bulk request
func (ap *accountsProcessor) SerializeAccounts(
	accounts map[string]*data.AccountInfo,
	areESDTAccounts bool,
) ([]*bytes.Buffer, error) {
	buffSlice := data.NewBufferSlice()
	for _, acc := range accounts {
		meta, serializedData, err := prepareSerializedAccount(acc, areESDTAccounts)
		if err != nil {
			return nil, err
		}

		err = buffSlice.PutData(meta, serializedData)
		if err != nil {
			return nil, err
		}
	}

	return buffSlice.Buffers(), nil
}

func prepareSerializedAccount(acc *data.AccountInfo, isESDT bool) ([]byte, []byte, error) {
	if (acc.Balance == "0" || acc.Balance == "") && isESDT {
		meta := prepareDeleteAccountInfo(acc, isESDT)
		return meta, nil, nil
	}

	return prepareSerializedAccountInfo(acc, isESDT)
}

func prepareDeleteAccountInfo(acct *data.AccountInfo, isESDT bool) []byte {
	id := acct.Address
	tokenName := adjustTokenIdentifierIfNeeded(acct.TokenName)
	if isESDT {
		hexEncodedNonce := converters.EncodeNonceToHex(acct.TokenNonce)
		id += fmt.Sprintf("-%s-%s", tokenName, hexEncodedNonce)
	}

	meta := []byte(fmt.Sprintf(`{ "delete" : { "_id" : "%s" } }%s`, id, "\n"))

	return meta
}

func adjustTokenIdentifierIfNeeded(identifier string) string {
	splitStr := strings.Split(identifier, esdtIdentifierSeparator)
	if len(splitStr) != 2 {
		return identifier
	}

	randomSequence := splitStr[1]
	if len(randomSequence) == esdtRandomSequenceLength {
		return identifier
	}

	if len(randomSequence) < esdtRandomSequenceLength {
		return identifier
	}

	return fmt.Sprintf("%s-%s", splitStr[0], randomSequence[:esdtRandomSequenceLength])
}

func prepareSerializedAccountInfo(
	account *data.AccountInfo,
	isESDTAccount bool,
) ([]byte, []byte, error) {
	id := account.Address
	if isESDTAccount {
		hexEncodedNonce := converters.EncodeNonceToHex(account.TokenNonce)
		id += fmt.Sprintf("-%s-%s", account.TokenName, hexEncodedNonce)
	}

	meta := []byte(fmt.Sprintf(`{ "index" : { "_id" : "%s" } }%s`, id, "\n"))
	serializedData, err := json.Marshal(account)
	if err != nil {
		return nil, nil, err
	}

	return meta, serializedData, nil
}

// SerializeAccountsHistory will serialize accounts history in a way that Elastic Search expects a bulk request
func (ap *accountsProcessor) SerializeAccountsHistory(accounts map[string]*data.AccountBalanceHistory) ([]*bytes.Buffer, error) {
	var err error

	buffSlice := data.NewBufferSlice()
	for _, acc := range accounts {
		meta, serializedData, errPrepareAcc := prepareSerializedAccountBalanceHistory(acc)
		if errPrepareAcc != nil {
			return nil, err
		}

		err = buffSlice.PutData(meta, serializedData)
		if err != nil {
			return nil, err
		}
	}

	return buffSlice.Buffers(), nil
}

func prepareSerializedAccountBalanceHistory(
	account *data.AccountBalanceHistory,
) ([]byte, []byte, error) {
	id := account.Address

	isESDT := account.Token != ""
	if isESDT {
		hexEncodedNonce := converters.EncodeNonceToHex(account.TokenNonce)
		id += fmt.Sprintf("-%s-%s", account.Token, hexEncodedNonce)
	}

	id += fmt.Sprintf("-%d", account.Timestamp)
	meta := []byte(fmt.Sprintf(`{ "index" : { "_id" : "%s" } }%s`, id, "\n"))

	serializedData, err := json.Marshal(account)
	if err != nil {
		return nil, nil, err
	}

	return meta, serializedData, nil
}
