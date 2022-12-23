package transactions

import (
	"encoding/hex"
	"testing"

	"github.com/ElrondNetwork/elastic-indexer-go/data"
	"github.com/stretchr/testify/require"
)

func TestAttachSCRsToTransactionsAndReturnSCRsWithoutTx(t *testing.T) {
	t.Parallel()

	scrsDataToTxs := newScrsDataToTransactions()

	txHash1 := []byte("txHash1")
	txHash2 := []byte("txHash2")
	tx1 := &data.Transaction{
		Hash:     hex.EncodeToString(txHash1),
		Nonce:    1,
		Sender:   "sender",
		Receiver: "receiver",
		GasLimit: 10000000,
		GasPrice: 1000000000,
		Data:     []byte("callSomething"),
		GasUsed:  5963500,
		Fee:      "128440000000000",
	}
	tx2 := &data.Transaction{}
	txs := map[string]*data.Transaction{
		string(txHash1): tx1,
		string(txHash2): tx2,
	}
	scrs := []*data.ScResult{
		{
			Nonce:          2,
			Sender:         "receiver",
			Receiver:       "sender",
			OriginalTxHash: hex.EncodeToString(txHash1),
			PrevTxHash:     hex.EncodeToString(txHash1),
			Value:          "40365000000000",
			Data:           []byte("@6f6b"),
		},
		{
			OriginalTxHash: "0102030405",
		},
	}

	scrsWithoutTx := scrsDataToTxs.attachSCRsToTransactionsAndReturnSCRsWithoutTx(txs, scrs)
	require.Len(t, scrsWithoutTx, 1)
	require.Len(t, tx1.SmartContractResults, 1)
	require.Equal(t, uint64(5963500), tx1.GasUsed)
	require.Equal(t, "128440000000000", tx1.Fee)

	require.Equal(t, scrsWithoutTx[0].OriginalTxHash, "0102030405")
}

func TestProcessTransactionsAfterSCRsWereAttached(t *testing.T) {
	t.Parallel()

	scrsDataToTxs := newScrsDataToTransactions()

	txHash1 := []byte("txHash1")
	txHash2 := []byte("txHash2")
	tx1 := &data.Transaction{
		Hash:     hex.EncodeToString(txHash1),
		Nonce:    1,
		Sender:   "sender",
		Receiver: "receiver",
		GasLimit: 10000000,
		GasPrice: 1000000000,
		Data:     []byte("callSomething"),
		SmartContractResults: []*data.ScResult{
			{
				ReturnMessage: "user error",
			},
		},
		GasUsed: 10000000,
		Fee:     "168805000000000",
	}
	tx2 := &data.Transaction{}
	txs := map[string]*data.Transaction{
		string(txHash1): tx1,
		string(txHash2): tx2,
	}

	scrsDataToTxs.processTransactionsAfterSCRsWereAttached(txs)
	require.Equal(t, "fail", tx1.Status)
	require.Equal(t, tx1.GasLimit, tx1.GasUsed)
	require.Equal(t, "168805000000000", tx1.Fee)
}

func TestIsESDTNFTTransferWithUserError(t *testing.T) {
	t.Parallel()

	require.False(t, isESDTNFTTransferOrMultiTransferWithError("ESDTNFTTransfer@45474c444d4558462d333766616239@06f5@045d2bd2629df0d2ea@0801120a00045d2bd2629df0d2ea226408f50d1a2000000000000000000500e809539d1d8febc54df4e6fe826fdc8ab6c88cf07ceb32003a3b00000007401c82df9c05a80000000000000407000000000000040f010000000009045d2bd2629df0d2ea0000000000000009045d2bd2629df0d2ea@636c61696d52657761726473"))
	require.False(t, isESDTNFTTransferOrMultiTransferWithError("ESDTTransfer@4d45582d623662623764@74b7e37e3c2efe5f11@"))
	require.False(t, isESDTNFTTransferOrMultiTransferWithError("ESDTNFTTransfer@45474c444d4558462d333766616239@070f@045d2bd2629df0d2ea@0801120a00045d2bd2629df0d2ea2264088f0e1a2000000000000000000500e809539d1d8febc54df4e6fe826fdc8ab6c88cf07ceb32003a3b000000074034d62af2b6930000000000000407000000000000040f010000000009045d2bd2629df0d2ea0000000000000009045d2bd2629df0d2ea@"))
	require.True(t, isESDTNFTTransferOrMultiTransferWithError("MultiESDTNFTTransfer@02@5745474c442d626434643739@00@38e62046fb1a0000@584d45582d666461333535@07@0801120c00048907e58284c28e898e2922520807120a4d45582d3435356335371a20000000000000000005007afb2c871d1647372fd53a9eb3e53e5a8ec9251cb05532003a1e0000000a4d45582d343535633537000000000000000000000000000008e8@657865637574696f6e206661696c6564"))
}
