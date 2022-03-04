TESTS_TO_RUN := $(shell go list ./... | grep -v integrationtests | grep -v mock)


test:
	@echo "  >  Running unit tests"
	go test -cover -race -coverprofile=coverage.txt -covermode=atomic -v ${TESTS_TO_RUN}

integration-tests:
	@echo " > Running integration tests"
	cd scripts && ./script.sh start
	go test -v ./integrationtests -tags integrationtests
	cd scripts && ./script.sh delete
	cd scripts && ./script.sh stop

start-cluster-with-kibana:
	@echo " > Starting Elasticsearch node and Kibana"
	docker-compose up -d

stop-cluster:
	docker-compose down
