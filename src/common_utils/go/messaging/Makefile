.PHONY: integration-test
integration-test:
	docker compose up -d
	sleep 10
	go test -v ./...
	docker compose down
