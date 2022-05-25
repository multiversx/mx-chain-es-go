package check

import (
	"github.com/ElrondNetwork/elastic-indexer-go/client/logging"
	"github.com/ElrondNetwork/elastic-indexer-go/converters"
	"github.com/ElrondNetwork/elastic-indexer-go/tools/accounts-balance-checker/pkg/config"
	"github.com/ElrondNetwork/elastic-indexer-go/tools/accounts-balance-checker/pkg/esclient"
	"github.com/ElrondNetwork/elastic-indexer-go/tools/accounts-balance-checker/pkg/rest"
	"github.com/ElrondNetwork/elrond-go-core/core/pubkeyConverter"
	"github.com/elastic/go-elasticsearch/v7"
)

func CreateBalanceChecker(cfg *config.Config, repair bool) (*balanceChecker, error) {
	esClient, err := esclient.NewElasticClient(elasticsearch.Config{
		Addresses: []string{cfg.Elasticsearch.URL},
		Username:  cfg.Elasticsearch.Username,
		Password:  cfg.Elasticsearch.Password,
		Logger:    &logging.CustomLogger{},
	})
	if err != nil {
		return nil, err
	}

	restClient, err := rest.NewRestClient(cfg.Proxy.URL)
	if err != nil {
		return nil, err
	}

	pubKeyConverter, err := pubkeyConverter.NewBech32PubkeyConverter(32, log)
	if err != nil {
		return nil, err
	}

	balanceToFloat, err := converters.NewBalanceConverter(18)
	if err != nil {
		return nil, err
	}

	return NewBalanceChecker(esClient, restClient, pubKeyConverter, balanceToFloat, repair)
}