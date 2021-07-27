package logsevents

import (
	"fmt"

	"github.com/ElrondNetwork/elastic-indexer-go/converters"
	"github.com/ElrondNetwork/elastic-indexer-go/data"
)

type pendingBalancesProc struct {
	pendingBalances map[string]*data.AccountInfo
}

func newPendingBalancesProcessor() *pendingBalancesProc {
	return &pendingBalancesProc{
		pendingBalances: make(map[string]*data.AccountInfo),
	}
}

func (pbp *pendingBalancesProc) addInfo(receiver string, token string, tokenNonce uint64, value string) {
	hexEncodedNonce := converters.EncodeNonceToHex(tokenNonce)
	key := fmt.Sprintf("%s_%s_%s_%s", pendingBalanceIdentifier, receiver, token, hexEncodedNonce)

	pbp.pendingBalances[key] = &data.AccountInfo{
		Address:         fmt.Sprintf("%s_%s", pendingBalanceIdentifier, receiver),
		Balance:         value,
		TokenName:       token,
		TokenIdentifier: converters.ComputeTokenIdentifier(token, tokenNonce),
		TokenNonce:      tokenNonce,
	}
}

func (pbp *pendingBalancesProc) getAll() map[string]*data.AccountInfo {
	return pbp.pendingBalances
}
