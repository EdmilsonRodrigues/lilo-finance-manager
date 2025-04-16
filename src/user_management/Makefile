.PHONY: lint
lint:
	ruff format
	ruff check --fix
	ruff format

.PHONY: mypy
mypy:
	mypy app/

.PHONY: e2e-test
e2e-test:
	python -m unittest discover -p 'test_*.py' -s test/e2e

.PHONY: integration-test
integration-test:
	echo "Running integration tests"
	# python -m unittest discover -p 'test_*.py' -s test/integration

.PHONY: unit-test
unit-test:
	python -m unittest discover -p 'test_*.py' -s test/unit

.PHONY: test
test: unit-test integration-test e2e-test


.PHONY: dev-test
dev-test: lint mypy test
