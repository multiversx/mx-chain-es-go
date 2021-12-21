package withKibana

// SCResults will hold the configuration for the scresults index
var SCResults = Object{
	"index_patterns": Array{
		"scresults-*",
	},
	"settings": Object{
		"number_of_shards":   3,
		"number_of_replicas": 0,
		"index": Object{
			"sort.field": Array{
				"timestamp",
			},
			"sort.order": Array{
				"desc",
			},
		},
		"opendistro.index_state_management.rollover_alias": "scresults",
	},
	"mappings": Object{
		"properties": Object{
			"timestamp": Object{
				"type":   "date",
				"format": "epoch_second",
			},
			"nonce": Object{
				"type": "unsigned_long",
			},
			"gasLimit": Object{
				"type": "unsigned_long",
			},
			"gasPrice": Object{
				"type": "unsigned_long",
			},
		},
	},
}
