.PHONY: lint
lint:
	ruff format
	ruff check --fix
	ruff format

.PHONY: mypy
mypy:
	mypy app/

.PHONY: test
test:
	python -m unittest discover -p 'test_*.py' -s test

.PHONY: integration-test
integration-test:
	echo "Running integration tests"
	# python -m unittest discover -p 'test_*.py' -s test/integration

.PHONY: unit-test
unit-test:
	python -m unittest discover -p 'test_*.py' -s test/unit

.PHONY: dev-test
dev-test: lint mypy test
