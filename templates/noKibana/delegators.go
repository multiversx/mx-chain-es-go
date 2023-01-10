package noKibana

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
			"timestamp": Object{
				"index":  "false",
				"type":   "date",
				"format": "epoch_second",
			},
			"unDelegateInfo": Object{
				"properties": Object{
					"id": Object{
						"index": "false",
						"type":  "keyword",
					},
					"value": Object{
						"index": "false",
						"type":  "keyword",
					},
					"valueNum": Object{
						"index": "false",
						"type":  "double",
					},
					"timestamp": Object{
						"index":  "false",
						"type":   "date",
						"format": "epoch_second",
					},
				},
			},
		},
	},
}
