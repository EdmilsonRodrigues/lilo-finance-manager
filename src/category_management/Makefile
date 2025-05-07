TAGS="-v"

.PHONY: e2e-test
e2e-test:
	GIN_MODE=release go test ${TAGS} -tags=e2e ./...

.PHONY: integration-test
integration-test:
	GIN_MODE=release go test ${TAGS} -tags=integration ./...

.PHONY: unit-test
unit-test:
	GIN_MODE=release go test ${TAGS} -tags=unit ./...


.PHONY: test
test: unit-test integration-test e2e-test


.PHONY: dev-test
dev-test: lint mypy test
