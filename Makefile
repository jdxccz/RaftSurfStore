.SILENT:

.PHONY: run-server
run-server:
	go run cmd/server/main.go -f serverConfig.json -m all

.PHONY: run-client
run-client: %:
	go run cmd/client/main.go -f clientConfig.json -b $(ARGS)

ARGS = $(filter-out $@,$(MAKECMDGOALS))