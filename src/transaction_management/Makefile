.PHONY: lint
lint:
	ruff format
	ruff check --fix
	ruff format

.PHONY: mypy
mypy:
	mypy setup/ transaction_manager/

.PHONY: e2e-test
e2e-test:
	# python -m unittest discover -p 'test_*.py' -s test/e2e -v

.PHONY: integration-test
integration-test:
	# python -m unittest discover -p 'test_*.py' -s test/integration -v

.PHONY: unit-test
unit-test:
	# python -m unittest discover -p 'test_*.py' -s test/unit -v

.PHONY: test
test: unit-test integration-test e2e-test


.PHONY: dev-test
dev-test: lint mypy test
