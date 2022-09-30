package dataindexer

import (
	"math/big"

	"github.com/ElrondNetwork/elastic-indexer-go/data"
	"github.com/ElrondNetwork/elastic-indexer-go/process/dataindexer/workItems"
	coreData "github.com/ElrondNetwork/elrond-go-core/data"
	"github.com/ElrondNetwork/elrond-go-core/data/block"
	"github.com/ElrondNetwork/elrond-go-core/data/outport"
)

// DispatcherHandler defines the interface for the dispatcher that will manage when items are saved in elasticsearch database
type DispatcherHandler interface {
	StartIndexData()
	Close() error
	Add(item workItems.WorkItemHandler)
	IsInterfaceNil() bool
}

// ElasticProcessor defines the interface for the elastic search indexer
type ElasticProcessor interface {
	SaveHeader(
		headerHash []byte,
		header coreData.HeaderHandler,
		signersIndexes []uint64,
		body *block.Body,
		notarizedHeadersHashes []string,
		gasConsumptionData outport.HeaderGasConsumption,
		txsSize int,
	) error
	RemoveHeader(header coreData.HeaderHandler) error
	RemoveMiniblocks(header coreData.HeaderHandler, body *block.Body) error
	RemoveTransactions(header coreData.HeaderHandler, body *block.Body) error
	RemoveAccountsESDT(headerTimestamp uint64, shardID uint32) error
	SaveMiniblocks(header coreData.HeaderHandler, body *block.Body) error
	SaveTransactions(body *block.Body, header coreData.HeaderHandler, pool *outport.Pool, coreAlteredAccounts map[string]*outport.AlteredAccount, isImportDB bool, numOfShards uint32) error
	SaveValidatorsRating(index string, validatorsRatingInfo []*data.ValidatorRatingInfo) error
	SaveRoundsInfo(infos []*data.RoundInfo) error
	SaveShardValidatorsPubKeys(shardID, epoch uint32, shardValidatorsPubKeys [][]byte) error
	SaveAccounts(blockTimestamp uint64, accounts []*data.Account, shardID uint32) error
	IsInterfaceNil() bool
}

// FeesProcessorHandler defines the interface for the transaction fees processor
type FeesProcessorHandler interface {
	ComputeGasUsedAndFeeBasedOnRefundValue(tx coreData.TransactionWithFeeHandler, refundValue *big.Int) (uint64, *big.Int)
	ComputeTxFeeBasedOnGasUsed(tx coreData.TransactionWithFeeHandler, gasUsed uint64) *big.Int
	ComputeGasLimit(tx coreData.TransactionWithFeeHandler) uint64
	IsInterfaceNil() bool
}

// ShardCoordinator defines what a shard state coordinator should hold
type ShardCoordinator interface {
	NumberOfShards() uint32
	ComputeId(address []byte) uint32
	SelfId() uint32
	SameShard(firstAddress, secondAddress []byte) bool
	CommunicationIdentifier(destShardID uint32) string
	IsInterfaceNil() bool
}

// Indexer is an interface for saving node specific data to other storage.
// This could be an elastic search index, a MySql database or any other external services.
type Indexer interface {
	SaveBlock(args *outport.ArgsSaveBlockData) error
	RevertIndexedBlock(header coreData.HeaderHandler, body coreData.BodyHandler) error
	SaveRoundsInfo(roundsInfos []*outport.RoundInfo) error
	SaveValidatorsPubKeys(validatorsPubKeys map[uint32][][]byte, epoch uint32) error
	SaveValidatorsRating(indexID string, infoRating []*outport.ValidatorRatingInfo) error
	SaveAccounts(blockTimestamp uint64, acc map[string]*outport.AlteredAccount, shardID uint32) error
	FinalizedBlock(headerHash []byte) error
	Close() error
	IsInterfaceNil() bool
	IsNilIndexer() bool
}

// BalanceConverter defines what a balance converter should be able to do
type BalanceConverter interface {
	ComputeBalanceAsFloat(balance *big.Int) float64
	ComputeESDTBalanceAsFloat(balance *big.Int) float64
	IsInterfaceNil() bool
}