package withKibana

// Delegators will hold the configuration for the delegators index
var Delegators = Object{
	"index_patterns": Array{
		"delegators-*",
	},
	"settings": Object{
		"number_of_shards":   3,
		"number_of_replicas": 0,
	},

	"mappings": Object{
		"properties": Object{
			"activeStake": Object{
				"type": "keyword",
			},
			"activeStakeNum": Object{
				"type": "double",
			},
			"address": Object{
				"type": "keyword",
			},
			"contract": Object{
				"type": "keyword",
			},
		},
	},
}
