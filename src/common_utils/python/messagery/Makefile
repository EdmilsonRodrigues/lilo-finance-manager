-PHONY: lint
lint:
	ruff format
	ruff check --fix
	ruff format

.PHONY: mypy
mypy:
	mypy src/

.PHONY: test
test:
	pytest

.PHONY: unit-test
unit-test:
	pytest test/test_kafka_unit.py

.PHONY: integration-test
integration-test:
	pytest test/test_kafka_integration.py

.PHONY: dev-test
dev-test: lint mypy test
