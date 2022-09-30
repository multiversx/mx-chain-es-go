package factory

import (
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/ElrondNetwork/elastic-indexer-go/client"
	"github.com/ElrondNetwork/elastic-indexer-go/client/logging"
	"github.com/ElrondNetwork/elastic-indexer-go/process/dataindexer"
	"github.com/ElrondNetwork/elastic-indexer-go/process/elasticproc/factory"
	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/core/check"
	"github.com/ElrondNetwork/elrond-go-core/hashing"
	"github.com/ElrondNetwork/elrond-go-core/marshal"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/elastic/go-elasticsearch/v7"
)

var log = logger.GetOrCreate("indexer/factory")

// ArgsIndexerFactory holds all dependencies required by the data indexer factory in order to create
// new instances
type ArgsIndexerFactory struct {
	Enabled                  bool
	UseKibana                bool
	IndexerCacheSize         int
	Denomination             int
	BulkRequestMaxSize       int
	Url                      string
	UserName                 string
	Password                 string
	TemplatesPath            string
	EnabledIndexes           []string
	Marshalizer              marshal.Marshalizer
	Hasher                   hashing.Hasher
	AddressPubkeyConverter   core.PubkeyConverter
	ValidatorPubkeyConverter core.PubkeyConverter
}

// NewIndexer will create a new instance of Indexer
func NewIndexer(args ArgsIndexerFactory) (dataindexer.Indexer, error) {
	err := checkDataIndexerParams(args)
	if err != nil {
		return nil, err
	}

	elasticProcessor, err := createElasticProcessor(args)
	if err != nil {
		return nil, err
	}

	dispatcher, err := dataindexer.NewDataDispatcher(args.IndexerCacheSize)
	if err != nil {
		return nil, err
	}

	dispatcher.StartIndexData()

	arguments := dataindexer.ArgDataIndexer{
		Marshalizer:      args.Marshalizer,
		ElasticProcessor: elasticProcessor,
		DataDispatcher:   dispatcher,
	}

	return dataindexer.NewDataIndexer(arguments)
}

func retryBackOff(attempt int) time.Duration {
	d := time.Duration(math.Exp2(float64(attempt))) * time.Second
	log.Debug("elastic: retry backoff", "attempt", attempt, "sleep duration", d)

	return d
}

func createElasticProcessor(args ArgsIndexerFactory) (dataindexer.ElasticProcessor, error) {
	databaseClient, err := client.NewElasticClient(elasticsearch.Config{
		Addresses:     []string{args.Url},
		Username:      args.UserName,
		Password:      args.Password,
		Logger:        &logging.CustomLogger{},
		RetryOnStatus: []int{http.StatusConflict},
		RetryBackoff:  retryBackOff,
	})
	if err != nil {
		return nil, err
	}

	argsElasticProcFac := factory.ArgElasticProcessorFactory{
		Marshalizer:              args.Marshalizer,
		Hasher:                   args.Hasher,
		AddressPubkeyConverter:   args.AddressPubkeyConverter,
		ValidatorPubkeyConverter: args.ValidatorPubkeyConverter,
		UseKibana:                args.UseKibana,
		DBClient:                 databaseClient,
		Denomination:             args.Denomination,
		EnabledIndexes:           args.EnabledIndexes,
		BulkRequestMaxSize:       args.BulkRequestMaxSize,
	}

	return factory.CreateElasticProcessor(argsElasticProcFac)
}

func checkDataIndexerParams(arguments ArgsIndexerFactory) error {
	if arguments.IndexerCacheSize < 0 {
		return dataindexer.ErrNegativeCacheSize
	}
	if check.IfNil(arguments.AddressPubkeyConverter) {
		return fmt.Errorf("%w when setting AddressPubkeyConverter in indexer", dataindexer.ErrNilPubkeyConverter)
	}
	if check.IfNil(arguments.ValidatorPubkeyConverter) {
		return fmt.Errorf("%w when setting ValidatorPubkeyConverter in indexer", dataindexer.ErrNilPubkeyConverter)
	}
	if arguments.Url == "" {
		return dataindexer.ErrNilUrl
	}
	if check.IfNil(arguments.Marshalizer) {
		return dataindexer.ErrNilMarshalizer
	}
	if check.IfNil(arguments.Hasher) {
		return dataindexer.ErrNilHasher
	}

	return nil
}