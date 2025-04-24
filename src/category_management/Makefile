.PHONY: e2e-test
e2e-test:
	go test -v -tags=e2e ./...

.PHONY: integration-test
integration-test:
	go test -v -tags=integration ./...

.PHONY: unit-test
unit-test:
	go test -v -tags=unit ./...


.PHONY: test
test: unit-test integration-test e2e-test


.PHONY: dev-test
dev-test: lint mypy test
