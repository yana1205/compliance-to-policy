.PHONY:
build:
	python -m build

.PHONY:
install:
	python -m pip install .

.PHONY:
install-dev:
	python -m pip install ".[dev]"

.PHONY:
uninstall:
	python -m pip uninstall compliance-to-policy

.PHONY:
format:
	python -m isort .
	python -m black .

.PHONY:
lint:
	python -m pylint ./c2p ./tests

.PHONY: docs
docs:
	python -m mkdocs build

.PHONY: gh-pages
 gh-pages:
	python -m mkdocs gh-deploy

# make test ARGS="-n 2 --dist loadscope --log-cli-level DEBUG" TARGET="tests/c2p/test_cli.py"
# TODO: -n 2 (pytest-xdist plugin) results in no logs displayed.
test: ARGS ?= 
test: TARGET ?= tests/
test: test-plugin
	@OUTPUT_PATH=/dev/null $(PYTHON) -m pytest $(ARGS) $(TARGET)

test-plugin: ARGS ?= 
test-plugin: TARGET ?= plugins_public/tests/
test-plugin:
	@OUTPUT_PATH=/dev/null $(PYTHON) -m pytest $(ARGS) $(TARGET)

.PHONY:
it:
	python samples_public/kyverno/compliance_to_policy.py
	python samples_public/kyverno/result_to_compliance.py
	python samples_public/ocm/compliance_to_policy.py
	python samples_public/ocm/result_to_compliance.py
	python samples_public/auditree/compliance_to_policy.py
	python samples_public/auditree/result_to_compliance.py

clean:
	@rm -rf build *.egg-info dist
	@find ./plugins -type d \( -name '*.egg-info' -o -name 'dist' \) | while read x; do echo $$x; rm -r $$x ; done 
	python -m pyclean -v .