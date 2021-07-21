package logsevents

import (
	"encoding/hex"
	"encoding/json"
	"math/big"
	"testing"

	"github.com/ElrondNetwork/elastic-indexer-go/data"
	"github.com/ElrondNetwork/elastic-indexer-go/mock"
	"github.com/ElrondNetwork/elastic-indexer-go/process/tags"
	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/data/transaction"
	"github.com/ElrondNetwork/elrond-vm-common/data/esdt"
	"github.com/stretchr/testify/require"
)

func TestNftsProcessor_processLogAndEventsNFTs(t *testing.T) {
	t.Parallel()

	esdtData := &esdt.ESDigitalToken{
		TokenMetaData: &esdt.MetaData{
			Creator: []byte("creator"),
		},
	}
	esdtDataBytes, _ := json.Marshal(esdtData)

	nonce := uint64(19)
	event := &transaction.Event{
		Address:    []byte("addr"),
		Identifier: []byte(core.BuiltInFunctionESDTNFTCreate),
		Topics:     [][]byte{[]byte("my-token"), big.NewInt(0).SetUint64(nonce).Bytes(), esdtDataBytes},
	}

	nftsProc := newNFTsProcessor(&mock.ShardCoordinatorMock{}, &mock.PubkeyConverterMock{}, &mock.MarshalizerMock{})

	altered := data.NewAlteredAccounts()

	tagsCount := tags.NewTagsCount()
	tokensCreateInfo := data.NewTokensInfo()
	nftsProc.processEvent(&argsProcessEvent{
		event:     event,
		accounts:  altered,
		tokens:    tokensCreateInfo,
		timestamp: 1000,
		tagsCount: tagsCount,
	})

	alteredAddr, ok := altered.Get("61646472")
	require.True(t, ok)
	require.Equal(t, &data.AlteredAccount{
		IsNFTOperation:  true,
		TokenIdentifier: "my-token",
		NFTNonce:        19,
	}, alteredAddr[0])

	require.Equal(t, &data.TokenInfo{
		Identifier: "my-token-13",
		Token:      "my-token",
		Timestamp:  1000,
		Issuer:     "",
		Data: &data.TokenMetaData{
			Creator: hex.EncodeToString([]byte("creator")),
		},
	}, tokensCreateInfo.GetAll()[0])

}

func TestNftsProcessor_processLogAndEventsNFTs_TransferNFT(t *testing.T) {
	t.Parallel()

	nonce := uint64(19)
	nftsProc := newNFTsProcessor(&mock.ShardCoordinatorMock{}, &mock.PubkeyConverterMock{}, &mock.MarshalizerMock{})

	events := &transaction.Event{
		Address:    []byte("addr"),
		Identifier: []byte(core.BuiltInFunctionESDTNFTTransfer),
		Topics:     [][]byte{[]byte("my-token"), big.NewInt(0).SetUint64(nonce).Bytes(), []byte("receiver")},
	}

	altered := data.NewAlteredAccounts()

	tagsCount := tags.NewTagsCount()
	nftsProc.processEvent(&argsProcessEvent{
		event:     events,
		accounts:  altered,
		timestamp: 10000,
		tagsCount: tagsCount,
	})

	alteredAddrSender, ok := altered.Get("61646472")
	require.True(t, ok)
	require.Equal(t, &data.AlteredAccount{
		IsNFTOperation:  true,
		TokenIdentifier: "my-token",
		NFTNonce:        19,
	}, alteredAddrSender[0])

	alteredAddrReceiver, ok := altered.Get("7265636569766572")
	require.True(t, ok)
	require.Equal(t, &data.AlteredAccount{
		IsNFTOperation:  true,
		TokenIdentifier: "my-token",
		NFTNonce:        19,
	}, alteredAddrReceiver[0])
}
