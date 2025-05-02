.PHONY: integration-test
integration-test:
	docker compose up -d
	sleep 20
	go test -v ./...
	docker compose down
