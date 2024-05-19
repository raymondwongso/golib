.PHONY: dev
dev:
	bin/dev.sh

.PHONY: format
format:
	bin/format.sh

.PHONY: lint
lint:
	bin/lint.sh

.PHONY: test
test:
	bin/test.sh
